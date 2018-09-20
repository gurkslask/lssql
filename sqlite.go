package main

import (
	"database/sql"
	"fmt"
)

type sqlite struct {
	db     *sql.DB
	dbtype string
}

//Returns string with all available tables from db
func (d sqlite) availableTables() (string, error) {
	if *debugp {
		fmt.Println("In printAvailableTables")
	}
	var result string
	rows, err := d.db.Query(`SELECT name FROM sqlite_master WHERE type = "table"`)
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

func (d sqlite) columnInfo(tablename *string) ([][]string, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(fmt.Sprintf("PRAGMA TABLE_INFO(%s)", tablename))
	if err != nil {
		fmt.Println("stmt")
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("rows")
		return nil, err
	}
	heads, err := getData(rows)
	if err != nil {
		return nil, err
	}
	return heads, nil
}
func (d sqlite) getdb() sql.DB {
	return *d.db
}
func (d sqlite) getdbtype() string {
	return d.dbtype
}
