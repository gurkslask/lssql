package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	// "os"
)

func main() {
	fmt.Println("vim-go")
	var path = flag.String("path", "./", "Path to .db file")
	var table = flag.String("table", "default", "Name of the table in db")
	flag.Parse()
	db, err := connectDB(path)
	if err != nil {
		fmt.Println(err)
	}
	printTable(db, table)

}

func connectDB(path *string) (*sql.DB, error) {
	fmt.Println("Connect to a database")
	db, err := sql.Open("sqlite3", *path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func printTable(db *sql.DB, tablename *string) (string, error) {
	fmt.Println("Table")
	return "", nil
}
