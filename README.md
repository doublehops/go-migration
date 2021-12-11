## Golang Database Migration Tool

This is a tool that is designed to automate database migrations. It is based on
the design of migrations of Yii2 framework.

### Usage
`<entry_script> -action <action> <-name|-number> XXXX`

Create new migration file:  
`./migrate -action create -name add_user_table`  

Migrate pending migrations:  
`./migrate -action up -number 1` // The number is optional. Will run all migrations if 
value for number is not included.  

Rollback past migrations:  
`./migrate -action down -number 1` // number is optional. Will only rollback one migration if no number value is included.  

The library keeps track of which migrations have been run by creating and populating a table named `migrations`.

### Getting Started
The library can be included into your application with the following example. A database connection and 
path to the migration files need to be passed in.

An example of usage can be seen in the file `cmd/main.go` which can be modified and added into your own project.

The script can be ran with `go run cmd/main.go -action create -name new_user_table`

Migration files will be saved to the location defined by the second parameter in call to `migrate.New()`. Two files will be created.
One for migration up and another for down. Once created, you should edit the files with the raw SQL queries. Multiple
queries should be separated with a string of `------------------`
on a line of its own.

### Additional

The migrations are run inside a transaction and will attempt to rollback all statements if one fails. However, some
MySQL statements force a commit, preventing this from working as intended with MySQL/MariaDB databases.

### Todo
- Add tests
- split the files up into separate packages.
