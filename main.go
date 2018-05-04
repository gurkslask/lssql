package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	// "os"
)

func main() {
	var path = flag.String("path", "/home/alex/data.db", "Path to .db file")
	var table = flag.String("table", "", "Name of the table in db")
	flag.Parse()
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

type foo struct {
	test string
}

func connectDB(path *string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", *path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func printTable(db *sql.DB, tablename *string) (string, error) {
	var result string
	fmt.Println("Table")
	rows, err := db.Query("select * from PEOPLE")
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var id, age int
		var name string
		err = rows.Scan(&id, &name, &age)
		if err != nil {
			return "", err
		}
		result += fmt.Sprint("+", id, "+", name, "+", age, "+\n")
	}
	return result, nil
}

func printTables(db *sql.DB) (string, error) {
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
