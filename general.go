package main

import (
	"database/sql"
	"errors"
	"fmt"
)

//Gets strings from unknown columns
func getData(rows *sql.Rows) ([][]string, error) {
	if *debugp {
		fmt.Println("in getData")
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
	if *debugp {
		fmt.Println("End of getData")
	}
	return rresult, nil
}

//Appends data with length padding to dest
func padString(data string, length int, dest *string) {
	*dest += fmt.Sprintf("%-[1]*s\t", length, data)
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

//Prints help monologe
func printHelp() {
	fmt.Println(`
	NAME
		lssql - List SQL contents
	SYNOPSIS
		lssql [FILE] [OPTION]...
	DESCRIPTION
		List contents of SQL databases.

		-table 
		List contents of this table, if omitted print available tables

		-limit 
		Number of lines to print

		-offset
		Offset from where to start printing

		-debug 
		Prints extra debug info

		-dbtype
		Choose between sqlite and postgres (Default: sqlite)

		-help
		Prints help dialog (this)

`)
}
func getDbSpecifics(dbType string) (*dsa, error) {

	//var databasep dsa
	databasep := new(dsa)
	switch dbType {
	case "sqlite":
		psqlite := new(sqlite)
		databasep.lister = psqlite
		databasep.dbtype = "sqlite3"
	case "postgres":
		ppostgres := new(postgres)
		databasep.lister = ppostgres
		databasep.dbtype = "postgres"
	default:
		e := errors.New(fmt.Sprintf("No type with the name %s supported", dbType))
		return databasep, e
	}

	return databasep, nil
}

//Connect to database and return a db
func connectDB(path *string, specifiedDb *dsa) error {
	var err error
	if *debugp {
		fmt.Println("In connect")
	}
	specifiedDb.db, err = sql.Open(specifiedDb.dbtype, *path)
	if err != nil {
		return err
	}
	return nil
}
