package lssql

import (
	"database/sql"
	"fmt"
)

type Postgres struct {
	db     *sql.DB
	dbtype string
}

func (d Postgres) AvailableTables(db *sql.DB) (string, error) {
	var result string
	var schemaname, tablename, tableowner, tablespace string
	rows, err := db.Query(`SELECT * FROM pg_catalog.pg_tables`)
	if err != nil {
		return "", err
	}
	data, err := GetData(rows)
	if err != nil {
		return "", err
	}
	for _, row := range data {
		schemaname = row[0]
		tablename = row[1]
		tableowner = row[2]
		tablespace = row[3]
		result += fmt.Sprintf("Schemaname: %s, Tablename: %s, Tableowner: %s, Tablespace: %s\n", schemaname, tablename, tableowner, tablespace)
	}
	return result, nil
}
func (d Postgres) ColumnInfo(tablename *string, db *sql.DB) ([]DBhead, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(fmt.Sprintf("SELECT * FROM %s WHERE false", *tablename))
	if err != nil {
		fmt.Println("stmt")
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("rows")
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columntypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	//var heads [][]string
	var heads []DBhead
	for index, name := range columns {
		fmt.Println(name, columntypes[index].ScanType())
		//heads = append(heads, []string{name, columntypes[index].ScanType().Name()})
		heads = append(heads, DBhead{
			Colname: name,
			Coltype: columntypes[index].ScanType().Name(),
		})
		//heads = [][]string{{fmt.Sprint(name, columntypes[index].ScanType())}}
	}
	if err != nil {
		return nil, err
	}
	return heads, nil
}

func (d Postgres) Statement() string { return "SELECT * FROM %s LIMIT $1 OFFSET $2 " }
