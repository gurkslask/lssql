package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gurkslask/lssql"
)

var debug bool = false
var debugp *bool = &debug

func main() {
	var err error

	var table = flag.String("table", "", "Name of the table in db")
	var limit = flag.Int("limit", 10, "Number of lines to print")
	var offset = flag.Int("offset", 1, "Offset from where to start printing")
	var debugf = flag.Bool("debug", false, "Prints extra debug info")
	var dbtype = flag.String("dbtype", "sqlite", `What kind of database?\nSupported databases: sqlite and postgres`)
	var config = flag.Bool("config", false, "Use a config file, if file doesnt exist, print default config")
	var path = flag.String("path", "", "Path to database or config file")
	flag.Parse()

	debugp = debugf
	if *debugp {
		fmt.Println("In main")
		fmt.Println(table)
	}

	if *config {
		ss := strings.Split(*path, ".")
		fileending := ss[len(ss)-1]
		var c *lssql.ConfigT
		if fileending == "yml" {
			var cy lssql.Config_yml
			var err error
			c, err = lssql.GetConfig(cy, path)
			if err != nil {
				panic(err)
			}

		}
		path = &c.Path
		offset = &c.Offset
		limit = &c.Limit
		table = &c.Table
		dbtype = &c.Dbtype

	}

	dbSpecifics, err := lssql.GetDbSpecifics(*dbtype)
	if err != nil {
		log.Fatal(err)
	}
	err = lssql.ConnectDB(path, dbSpecifics)
	defer dbSpecifics.DB.Close()
	if err != nil {
		fmt.Println(err)
	}
	if *table == "" {
		tables, err := dbSpecifics.Lister.AvailableTables(dbSpecifics.DB)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("No table provided, listed tables: \n", tables)
	} else {
		result, err := printTable(*dbSpecifics, table, limit, offset)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	}

}

func printTable(d lssql.DBdialect, tablename *string, limit, offset *int) (string, error) {
	// Get data with queries and print it nicely with padding
	if *debugp {
		fmt.Println("*************In printTable")
		fmt.Printf("In Printtable with tablename: %s\n", *tablename)

	}
	fmt.Printf("Table %s\n\n", *tablename)
	stmt, err := d.DB.Prepare(fmt.Sprintf(d.Lister.Statement(), *tablename))
	if err != nil {
		fmt.Println("Statement error")
		return "", err
	}
	rows, err := stmt.Query(*limit, *offset)
	if err != nil {
		return "", err
	}
	data, err := lssql.GetData(rows)
	if err != nil {
		return "", err
	}
	if data == nil {
		err := errors.New("No data from the query")
		return "", err
	}
	heads, err := d.Lister.ColumnInfo(tablename, d.DB)
	if err != nil {
		log.Fatal("heads ", err)
	}

	if *debugp {
		fmt.Println("This is the data:")
		fmt.Println(data)
		fmt.Println("This is the heads:")
		fmt.Println(heads)
	}
	t := [][]string{{heads[0].Colname, heads[0].Coltype}}
	columnlengths := lssql.MaxColumnLength(data, t)

	resultstring := ""
	for i := 0; i < len(data[0]); i++ {
		lssql.PadString(heads[i].Colname, columnlengths[i], &resultstring)
		//resultstring += "\t"
	}
	resultstring += "\n"
	for i := 0; i < len(data[0]); i++ {
		lssql.PadString(heads[i].Coltype, columnlengths[i], &resultstring)
	}

	resultstring += "\n"
	resultstring += "\n"
	for _, row := range data {
		for i, col := range row {
			lssql.PadString(col, columnlengths[i], &resultstring)
		}
		resultstring += "\n"
	}

	if *debugp {
		fmt.Printf("Out Printtable with tablename: %s\n", *tablename)
	}
	return resultstring, nil
}
