package migrations

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/doublehops/go-migration/helpers"
)

// createMigration will copy template file into new fil
func (h *Handle) createMigration(name string) error {
	currentTime := time.Now()
	curTime := fmt.Sprintf(currentTime.Format("20060102_150405_"))
	name = curTime + name + ".json"
	path := h.path + "/" + name

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

	os.Stderr.WriteString("Migration file created: " + name + "\n")

	return nil
}

// migrateUp will run new migration/s.
func (h Handle) migrateUp(args *Args) error {
	var err error
	var pendingFiles []string

	pendingFiles, err = h.getPendingMigrationFiles()
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", pendingFiles)

	migrationFiles, err := h.parseMigrations(pendingFiles)
	if err != nil {
		return err
	}

	i := 0

	for _, file := range migrationFiles {
		err = h.processFileUp(file)
		if err != nil {
			return err
		}

		i++ // Only perform number of migrations (files) equal to that supplied in argument
		if i == args.number {
			break
		}
	}

	return nil
}

// processFile will process the queries in the given file. It will attempt to rollback when there is
// and error with one of the queries.
func (h *Handle) processFileUp(file File) error {

	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction. %w", err)
	}
	defer tx.Rollback()

	helpers.PrintMsg(fmt.Sprintf("Migrating queries from: %s", file.Filename))
	for _, q := range file.Queries.Up {
		_, err = tx.Exec(q)
		if err != nil {
			return fmt.Errorf("\nthere was an error executing query. File: %s; query; %s; Error: %s", file.Filename, q, err)
		}
	}
	_, err = tx.Exec(InsertMigrationRecordIntoTableSQL, file.Filename)
	if err != nil {
		return fmt.Errorf("unable to update migration table with newly ran migration record. %w", err)
	}
	helpers.PrintMsg(" - Success\n")

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("there was an error committing transaction. %w", err)
	}

	return nil
}
