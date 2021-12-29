package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/doublehops/go-migration"

	_ "github.com/go-sql-driver/mysql"
)

/*
  This file just serves as an example of how you would add this library to your project.
 */

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", "dev", "pass12", "127.0.0.1", "cw")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error creating database connection. %s", err)
	}

	action := getAction()
	action.Path = "/home/b/workspace/cryptowatcher-2.0/migrations"
	action.DB = db

	err = action.Migrate()
	if err != nil {
		fmt.Printf("error running migrations. %s", err)
	}
}

// getAction will read the flags from the command and populate an Args struct. The available flags are `action`,
// `number` and `name`.
func getAction() *go_migration.Action {
	var args go_migration.Action

	act := flag.String("action", "", "The intended action")
	number := flag.Int("number", 0, "The number of migrations to run")
	name := flag.String("name", "", "The name of the new migration")
	flag.Parse()

	args.Action = *act
	args.Number = *number
	args.Name = *name

	if found := args.IsValidAction(args.Action); !found {
		args.PrintHelp()
	}

	if args.Action == "create" && args.Name == "" {
		args.PrintHelp()
	}

	if args.Action == "up" && args.Number == 0 {
		args.Number = 9999 // run them all if none defined.
	}

	if args.Action == "down" && args.Number == 0 {
		args.Number = 1 // run just one if none defined.
	}

	return &args
}
