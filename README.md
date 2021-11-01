## Golang Database Migration Tool

This is a tool that is designed to automate database migrations. It is based on
the migrations of Yii2 framework.

### Usage
`<entry_script> <action> <name|number>`

Create new migration file:  
`./migrate create add_user_table`  

Migrate pending migrations:  
`./migrate up 1` // The number is optional. Will run all migrations if value is not included.  

Rollback past migrations:  
`./migrate down 1` // number is optional. Will only rollback one migration if no value is included.  

The library keeps track of which migrations have been run by creating and populating a table named `migrations`.

### Getting Started
The library can be included into your application with the following example. A database connection and 
path to the migration files need to be passed in.
```golang
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	migrate "github.com/doublehops/go-migration"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", "username", "password", "127.0.0.1", "dbname")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error creating database connection. %s", err)
	}

	m, err := migrate.New(db, "/path/to/migrations")
	if err != nil {
		fmt.Printf("error creating migration instance. %s", err)
	}

	err = m.Migrate()
	if err != nil {
		fmt.Printf("error running migrations. %s", err)
	}
}
```
If this script was saved as `migrate.go`, you would run it with `go run migrate.go create new_user_table`

### Additional
Unfortunately JSON doesn't allow carriage returns inside the string so each SQL statement needs to be 
on a single line, hurting readability. Suggestions about how to get around this are welcome.

The migrations are run inside a transaction in attempt to rollback all statements if one fails. However, some
MySQL statements force a commit, preventing this from working as intended.

### Todo
- Add tests

