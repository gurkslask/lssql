package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var debug bool = false
var debugp *bool = &debug

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No arguments supplied")
		os.Exit(1)
	}

	//var path = flag.String("path", "", "Path to db")
	var table = flag.String("table", "", "Name of the table in db")
	var limit = flag.Int("limit", 10, "Number of lines to print")
	var offset = flag.Int("offset", 1, "Offset from where to start printing")
	var debugf = flag.Bool("debug", false, "Prints extra debug info")
	flag.CommandLine.Parse(os.Args[2:])
	fmt.Println(os.Args)
	path := &os.Args[1]
	p := string("-table2")
	path = &p
	f, err := os.Stat(*path)
	fmt.Println(f)
	fmt.Println(err)
	fmt.Println(*path)
	if os.IsNotExist(err) {
		fmt.Println("Supplied path does not exist")
		os.Exit(1)
	}
	fmt.Println(err)
	debugp = debugf

	if *debugp {
		fmt.Println("In main")
		fmt.Println(table)
	}

	db, err := connectDB(path)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}

	if *path == "" {
		err = errors.New("No path specified")
		if err != nil {
			log.Fatal(err)
		}
	}
	if *table == "" {
		tables, err := printTables(db)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("No table provided, listed tables: \n", tables)
	} else {
		result, err := printTable(db, table, limit, offset)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	}

}

//Connect to database and return a db
func connectDB(path *string) (*sql.DB, error) {
	if *debugp {
		fmt.Println("In connect")
	}
	db, err := sql.Open("sqlite3", *path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

//Gets strings from unknown columns
func getData(rows *sql.Rows) ([][]string, error) {
	/*if *debugp {
		fmt.Printf("In getdata with query %s\n", query)
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}*/
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

	/* dest is where Scan puts data, make dest contain pointers to rawResult that is []byte */
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
		/* If rawresult isnt nil, make them string */
		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}
		//Initialize and copy to a [][]string that contains all the rows
		rresult = append(rresult, make([]string, len(result)))
		copy(rresult[rownumber], result)
		rownumber += 1
	}
	/*if *debugp {
		fmt.Printf("out getdata with query %s\n", query)
	}*/
	return rresult, nil
}

func printTable(db *sql.DB, tablename *string, limit, offset *int) (string, error) {
	// Get data with queries and print it nicely with padding
	if *debugp {
		fmt.Println("*************In printTable")
		fmt.Println("In Printtable with tablename: %s", *tablename)

	}
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	fmt.Printf("Table %s\n\n", *tablename)
	//fmt.Print("PRAGMA TABLE_INFO(%s)", *tablename)
	stmt, err := tx.Prepare(fmt.Sprintf("PRAGMA TABLE_INFO(%s)", *tablename))
	if err != nil {
		fmt.Println("stmt")
		return "", err
	}
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("rows")
		return "", err
	}
	heads, err := getData(rows)
	if err != nil {
		return "", err
	}
	fmt.Printf("Table %s\n\n", *tablename)
	//data, err := getData(db, fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d ", *tablename, *limit, *offset))
	//stmt, err = db.Prepare("SELECT * FROM ? LIMIT (?) OFFSET (?) ")
	stmt, err = db.Prepare(fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ? ", *tablename))
	if err != nil {
		return "", err
	}
	rows, err = stmt.Query(*limit, *offset)
	if err != nil {
		return "", err
	}
	data, err := getData(rows)
	if err != nil {
		return "", err
	}
	if data == nil {
		err := errors.New("No data from the query")
		return "", err

	}
	if *debugp {
		fmt.Println("This is the data:")
		fmt.Println(data)
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
func printTables(db *sql.DB) (string, error) {
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
