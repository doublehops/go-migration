package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/doublehops/go-migration/action"
)

type Handle struct {
	db     *sql.DB
	path   string
	action *action.Action
}

type TableList []Table

type Table struct {
	Name string
}

func New(db *sql.DB, path string) (*Handle, error) {

	var handle *Handle
	args, err := getArguments()
	if err != nil {
		return handle, err
	}

	a := &action.Action{
		Action: args.Action,
		Number: args.Number,
		Name: args.Name,
		DB: db,
	}

	return &Handle{
		db:   db,
		path: path,
		action: a,
	}, nil
}

func (h *Handle) Migrate() error {

	err := h.ensureMigrationsTableExists()
	if err != nil {
		return err
	}

	if h.action.Action == "create" {
		err = h.action.CreateMigration(h.path)
		return err
	}

	pendingFiles, err := h.getPendingMigrationFiles()
	if err != nil {
		return err
	}
	migrationFiles, err := h.parseMigrations(pendingFiles)
	if err != nil {
		return err
	}

	if h.action.Action == "up" {
		if err = h.action.MigrateUp(migrationFiles); err != nil {
			return err
		}
	}

	return nil
}

// getArguments will read the arguments from the command and populate an Args struct. Possible options for arg 1 is `create`,
// `up` and `down`. For create, the second arg is the migration name. For up/down, the second argument is the number of
// migrations to perform.
func getArguments() (*action.Action, error) {
	var args action.Action

	argList := os.Args[1:]

	if len(argList) < 1 {
		printHelp()
	}

	possibleArgs := []string{
		"create",
		"up",
		"down",
	}

	args.Action = argList[0]
	if found := sliceContains(args.Action, possibleArgs); !found {
		printHelp()
	}

	if args.Action == "create" {
		if len(argList) < 2 {
			printHelp()
		}
		args.Name = argList[1]

		return &args, nil
	}

	if len(argList) > 1 {
		number, err := strconv.Atoi(argList[1])
		if err != nil {
			printHelp()
		}

		if err != nil {
			return &args, fmt.Errorf("unable to convert second argument to int. %s", err)
		}
		args.Number = number
	} else {
		args.Number = 0
	}

	return &args, nil
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

// ensureMigrationsTableExists to create table to track migrations.
func (h *Handle) ensureMigrationsTableExists() error {
	var tableList TableList
	rows, err := h.db.Query(action.CheckMigrationsTableExistsSQL)
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
