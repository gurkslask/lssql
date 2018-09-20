package main

import (
	"database/sql"
	"fmt"
)

type postgres struct {
	db     *sql.DB
	dbtype string
}

func (d postgres) availableTables() (string, error)             { return "", nil }
func (d postgres) columnInfo(tablename *string) (string, error) { return "", nil }

/*specifiedDb.name = "postgres"
specifiedDb.columnInfoQuery = fmt.Sprintf("SELECT * FROM %s WHERE false", *table)
specifiedDb.availableTablesQuery = `select * from pg_catalog.pg_tables`*/
