package migrations

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type Args struct {
	action string
	number int
	name   string
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
		err = h.createMigration(args.name)
		if err != nil {
			return err
		}
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

		if len(os.Args[:1]) < 2 {
			printHelp()
		}
		args.name = os.Args[1]

		return args, nil
	}

	if len(argList) > 1 {
		if _, err := strconv.Atoi(argList[1]); err == nil {
			printHelp()
		}
		number, err := strconv.Atoi(argList[1])
		if err != nil {
			return args, fmt.Errorf("unable to convert second argument to int. %s", err)
		}
		args.number = number
	} else {
		args.number = 0
	}

	return args, nil
}

// createMigration will copy template file into new fil
func (h *Handle) createMigration(name string) error {
	currentTime := time.Now()
	curTime := fmt.Sprintf(currentTime.Format("20060102_150405_"))
	name = curTime + name + ".json"
	path := h.path + "/" + name

	//// @todo: use relative path - https://forum.golangbridge.org/t/how-to-get-relative-path-from-runtime-caller/15690/5
	//template := pwd+"/"+templatePath
	template := `{
  "up": [
    "CREATE TABLE 'test' ( name VARCHAR(255))"
  ],
  "down": [
    "DROP TABLE 'test'"
  ]
}`

	err := ioutil.WriteFile(path, []byte(template), 0644)
	if err != nil {
		return fmt.Errorf("unable to write template file. %s", err)
	}

	os.Stderr.WriteString("Migration file created: " + path + "\n")

	return nil
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
