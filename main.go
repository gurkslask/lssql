package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var debug bool = false
var debugp *bool = &debug

func main() {
	var path = flag.String("path", "", "Path to .db file")
	var table = flag.String("table", "", "Name of the table in db")
	var debugf = flag.Bool("debug", false, "Prints extra debug info")
	flag.Parse()
	debugp = debugf

	if *debugp {
		fmt.Println("In main")
	}

	db, err := connectDB(path)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}

	if *table == "" {
		tables, err := printTables(db)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("No table provided, listed tables: \n", tables)
	} else {
		result, err := printTable(db, table)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	}

}

//Connect to database and return a db
func connectDB(path *string) (*sqlx.DB, error) {
	if *debugp {
		fmt.Println("In connect")
	}
	db, err := sqlx.Open("sqlite3", *path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

//Gets strings from unknown columns
func getData(db *sqlx.DB, query string) ([][]string, error) {
	if *debugp {
		fmt.Printf("In getdata with query %s\n", query)
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	if *debugp {
		fmt.Println("Query done")
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	rawResult := make([][]byte, len(columns))
	result := make([]string, len(columns))
	var rresult [][]string
	rownumber := 0

	dest := make([]interface{}, len(columns))
	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}
		//Initialize and copy
		rresult = append(rresult, make([]string, len(result)))
		copy(rresult[rownumber], result)
		rownumber += 1
	}
	if *debugp {
		fmt.Printf("out getdata with query %s\n", query)
	}
	return rresult, nil
}

// Get data with queries and print it nicely with padding
func printTable(db *sqlx.DB, tablename *string) (string, error) {
	if *debugp {
		fmt.Println("*************In printTable")
		fmt.Println("In Printtable with tablename: %s", *tablename)

	}
	heads, err := getData(db, fmt.Sprintf("PRAGMA table_info(%v)", *tablename))
	if err != nil {
		return "", err
	}

	fmt.Printf("Table %s\n\n", *tablename)

	data, err := getData(db, fmt.Sprintf("select * from %v", *tablename))
	if err != nil {
		return "", err
	}
	columnlengths := maxColumnLength(data, heads)

	resultstring := ""
	for i := 0; i < len(data[0]); i++ {
		padString(heads[i][1], columnlengths[i], &resultstring)
		//resultstring += "\t"
	}
	resultstring += "\n"
	for i := 0; i < len(data[0]); i++ {
		////resultstring += heads[i][2]
		padString(heads[i][2], columnlengths[i], &resultstring)
		//resultstring += "\t"
	}

	resultstring += "\n"
	resultstring += "\n"
	for _, row := range data {
		for i, col := range row {
			//resultstring += col
			//resultstring += "\t"
			padString(col, columnlengths[i], &resultstring)
		}
		resultstring += "\n"
	}

	if *debugp {
		fmt.Println("Out Printtable with tablename: %s", *tablename)
	}
	return resultstring, nil
}

//Appends data with length padding to dest
func padString(data string, length int, dest *string) {
	*dest += fmt.Sprintf("%-[1]*s\t", length, data)
}

//If no tables are found or requested, print available tables
func printTables(db *sqlx.DB) (string, error) {
	if *debugp {
		fmt.Println("In printTables")
	}
	var result string
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type = "table"`)
	if err != nil {
		return "", err
	}
	for rows.Next() {

		var name string
		err = rows.Scan(&name)
		if err != nil {
			return "", err
		}
		result += fmt.Sprint(name, "\n")
	}
	return result, nil
}

//Gets max column length used for displaying data
func maxColumnLength(datain ...[][]string) []int {
	resultLength := 0
	for _, data := range datain {
		if len(data[0]) > resultLength {
			resultLength = len(data[0])
		}
	}
	result := make([]int, resultLength)
	for _, data := range datain {
		for i, _ := range data {
			for j, col := range data[i] {
				if result[j] < len(col) {
					//If length of current row is bigger, update data
					result[j] = len(col) + 5
				}

			}
		}
	}
	return result

}
