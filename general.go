package lssql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
)

type dblister interface {
	ColumnInfo(*string, *sql.DB) ([]DBhead, error)
	AvailableTables(*sql.DB) (string, error)
	Statement() string
}
type DBdialect struct {
	DBtype string
	DB     *sql.DB
	Lister dblister
}

type DBhead struct {
	Colname string
	Coltype string
}

//Gets strings from unknown columns
func GetData(rows *sql.Rows) ([][]string, error) {
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
	return rresult, nil
}

//Appends data with length padding to dest
func PadString(data string, length int, dest *string) {
	*dest += fmt.Sprintf("%-[1]*s\t", length, data)
}

//Gets max column length used for displaying data
func MaxColumnLength(datain ...[][]string) []int {
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
func PrintHelp() {
	fmt.Println(` NAME
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
		Prints help dialog (this) `)
}
func GetDbSpecifics(dbType string) (*DBdialect, error) {
	databasep := new(DBdialect)
	switch dbType {
	case "sqlite":
		psqlite := new(Sqlite)
		databasep.Lister = psqlite
		databasep.DBtype = "sqlite3"
	case "postgres":
		ppostgres := new(Postgres)
		databasep.Lister = ppostgres
		databasep.DBtype = "postgres"
	default:
		e := errors.New(fmt.Sprintf("No type with the name %s supported", dbType))
		return databasep, e
	}

	return databasep, nil
}

//Connect to database and return a db
func ConnectDB(path *string, specifiedDb *DBdialect) error {
	var err error
	specifiedDb.DB, err = sql.Open(specifiedDb.DBtype, *path)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig(config Config, path *string) (*ConfigT, error) {
	b, err := ioutil.ReadFile(*path)
	if err != nil {
		return nil, err
	}
	c, err := config.ReadConfig(b)
	if err != nil {
		return nil, err
	}
	return c, nil
}
