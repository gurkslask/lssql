package main

import (
	//"database/sql"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	//"strings"
	// "os"
	//"strconv"
)

var debug bool = true

func main() {
	if debug {
		fmt.Println("In main")
	}
	var path = flag.String("path", "/home/alex/data.db", "Path to .db file")
	var table = flag.String("table", "", "Name of the table in db")
	flag.Parse()
	*table = "TEST"
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

func connectDB(path *string) (*sqlx.DB, error) {
	if debug {
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
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	rawResult := make([][]byte, len(columns))
	result := make([]string, len(columns))
	rresult := make([][]string, len(columns))
	// var rresult [][]string
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
		rresult[rownumber] = make([]string, len(result))
		copy(rresult[rownumber], result)
		rownumber += 1
	}
	return rresult, nil
}

func printTable(db *sqlx.DB, tablename *string) (string, error) {
	if debug {
		fmt.Println("*************In printTable")
	}

	heads, err := getData(db, "PRAGMA table_info(TEST)")
	if err != nil {
		return "", err
	}

	fmt.Println("Table")

	data, err := getData(db, "select * from TEST")
	if err != nil {
		return "", err
	}

	resultstring := ""
	for i := 0; i < len(data[0]); i++ {
		resultstring += heads[i][1]
		resultstring += "\t"
	}
	resultstring += "\n"
	for i := 0; i < len(data[0]); i++ {
		resultstring += heads[i][2]
		resultstring += "\t"
	}
	resultstring += "\n"
	for _, row := range data {
		for _, col := range row {
			resultstring += col
			resultstring += "\t"
		}
		resultstring += "\n"
	}

	return resultstring, nil
}

func printTables(db *sqlx.DB) (string, error) {
	if debug {
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
