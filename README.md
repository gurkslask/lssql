## lssql
#Terminal sqlite browser

Use lssql to see what your sql tables contains, just like ls for POSIX. If you dont know the tablename lssql will print available tables. 
lssql is great to use while developing your application with SQL, and you need to see a quick glance at your database, just fire up lssql. If you know beforehand what database and table you will want to see, make a .yml config file as a shortcut.

The lssql must always have a path after the command, either to a config file or a database file/server

Available flags
* * *
* -table    table to show, if omitted print available tables
* -debug    prints debug information
* -limit    Limits how many values should be printed from the database
* -offset   Offsets the values that are to be printed
* -dbtype   Type of database as a string (postgres, sqlite)
* -config   This is a config file

Things that would be nice to have
* * *
* -Tail -f like behaviour, as in lssql continously queries the database and prints the changes that happen
* -Support for more database types (MSSQL, MySQL, MariaDB etc.)
