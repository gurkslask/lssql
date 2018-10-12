package lssql

import (
	"database/sql"
	"fmt"
)

type Sqlite struct {
}

//Returns string with all available tables from db
func (d Sqlite) AvailableTables(db *sql.DB) (string, error) {
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

func (d Sqlite) ColumnInfo(tablename *string, db *sql.DB) ([]DBhead, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(fmt.Sprintf("PRAGMA TABLE_INFO(%s)", *tablename))
	if err != nil {
		fmt.Println("stmt")
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("rows")
		return nil, err
	}
	heads, err := GetData(rows)
	if err != nil {
		return nil, err
	}
	var t []DBhead
	for i, _ := range heads {
		t = append(t, DBhead{
			Colname: heads[i][1],
			Coltype: heads[i][2],
		})
	}
	_ = t
	return t, nil
}

func (d Sqlite) Statement() string { return "SELECT * FROM %s LIMIT ? OFFSET ? " }
