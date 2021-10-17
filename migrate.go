package migrations

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Args struct {
	action string
	number int
}

type Handle struct {
	db *sql.DB
	path string
}

type TableList []Table

type Table struct {
	Name string
}

func New(db *sql.DB, path string) *Handle {

	return &Handle{
		db: db,
		path: path,
	}
}

func (h *Handle) Migrate() error {
	args := getArguments()
	err := h.ensureMigrationsTable()
	if err != nil {
		fmt.Println(">>>>>> Failing here <<<<<<<")
		return err
	}

	lastRanMigration, err := h.getLatestRanMigration()
	if err != nil {
		return err
	}
	allFiles, err := h.listMigrationFiles()
	if err != nil {
		return err
	}

	switch args.action {
	case "create":
		fmt.Println("create")
	}

	migrationsNotRun := h.getMigrationsNotRun(allFiles, lastRanMigration)
	fmt.Printf("%v\n", migrationsNotRun)

	direction := flag.String("direction", "up", "Direction to migrate. `up` (default) or `down`")
	number := flag.Int("number", 0, "Number of migration files to process. Defaults to 0 (all pending for up, or 1 for down)")
	flag.Parse()

	if *direction == "up" {
		err = h.up(number, allFiles)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h Handle) up(number *int, allFiles []string) error {

	for _, file := range allFiles {
		fmt.Println(file)
	}

	return nil
}

// getArguments will read the arguments from the command and populate an Args struct.
func getArguments() *Args{
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
		return args
	}

	if len(argList) > 1 {
		num, err := strconv.Atoi(argList[1])
		if err != nil {
			printHelp()
		}
		args.number = num
	} else {
		args.number = 9999 // Run all migrations
	}

	return args
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
	fmt.Printf("%v\n", tableList)

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
