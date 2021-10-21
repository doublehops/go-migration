package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

type Args struct {
	action string
	number int
	name   string
}

type File struct {
	Filename string
	Queries  *Queries
}

type Queries struct {
	Up   []string `json:"up"`
	Down []string `json:"down"`
}

type Handle struct {
	db   *sql.DB
	path string
}

type TableList []Table

type Table struct {
	Name string
}

func New(db *sql.DB, path string) *Handle {

	return &Handle{
		db:   db,
		path: path,
	}
}

func (h *Handle) Migrate() error {
	args, err := getArguments()
	if err != nil {
		return err
	}

	err = h.ensureMigrationsTable()
	if err != nil {
		return err
	}

	switch args.action {
	case "create":
		err = h.createMigration(args.name)
	case "up":
		err = h.migrateUp(args)
	}
	if err != nil {
		return err
	}

	return nil
}

// getArguments will read the arguments from the command and populate an Args struct. Possible options for arg 1 is `create`,
// `up` and `down`. For create, the second arg is the migration name. For up/down, the second argument is the number of
// migrations to perform.
func getArguments() (*Args, error) {
	args := &Args{}

	argList := os.Args[1:]

	if len(argList) < 1 {
		printHelp()
	}

	possibleArgs := []string{
		"create",
		"up",
		"down",
	}

	args.action = argList[0]
	if found := sliceContains(args.action, possibleArgs); !found {
		printHelp()
	}

	if args.action == "create" {
		if len(argList) < 2 {
			printHelp()
		}
		args.name = argList[1]

		return args, nil
	}

	if len(argList) > 1 {
		number, err := strconv.Atoi(argList[1])
		if err != nil {
			printHelp()
		}

		if err != nil {
			return args, fmt.Errorf("unable to convert second argument to int. %s", err)
		}
		args.number = number
	} else {
		args.number = 0
	}

	return args, nil
}

func sliceContains(key string, slice []string) bool {
	for _, item := range slice {
		if item == key {
			return true
		}
	}

	return false
}

func printHelp() {
	var helpMsg = `
Usage: <<< show help here >>>
`
	os.Stderr.WriteString(helpMsg)
	os.Exit(1)
}

func (h *Handle) ensureMigrationsTable() error {
	var tableList TableList
	rows, err := h.db.Query(CheckMigrationsTableExistsSQL)
	if err != nil {
		return fmt.Errorf("++++>>>> Error: %w", err)
	}

	for rows.Next() {
		var t Table
		if err = rows.Scan(&t.Name); err != nil {
			return err
		}
		tableList = append(tableList, t)
	}

	for _, tbl := range tableList {
		if tbl.Name == "migrations" { // Already exists
			return nil
		}
	}

	err = h.addMigrationTable()
	if err != nil {
		return err
	}

	return nil
}
