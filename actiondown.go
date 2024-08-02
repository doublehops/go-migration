package go_migration

import (
	"fmt"
	"strings"

	"github.com/doublehops/go-migration/helpers"
)

// MigrateDown will rollback migration/s.
func (a *Action) MigrateDown(migrationFiles []File) error {
	var err error

	for _, file := range migrationFiles {
		err = a.processFileDown(file)
		if err != nil {
			return err
		}
	}

	return nil
}

// processFileDown will process the down queries in the given file. It will attempt to rollback when there is
// and error in one of the queries.
func (a *Action) processFileDown(file File) error {

	tx, err := a.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction. %w", err)
	}
	defer tx.Rollback() // nolint

	helpers.PrintMsg(fmt.Sprintf("Migrating down queries from: %s", file.Filename))
	for _, q := range file.Queries {
		_, err = tx.Exec(q)
		if err != nil {
			return fmt.Errorf("\nthere was an error executing query. File: %s; Error: %s", file.Filename, err)
		}
	}
	filename := strings.Replace(file.Filename, ".down.sql", ".up.sql", 1)
	_, err = tx.Exec(RemoveMigrationRecordFromTableSQL, filename)
	if err != nil {
		return fmt.Errorf("unable to remove from migration table with newly ran migration record. %w", err)
	}
	helpers.PrintMsg(" - Success\n")

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("there was an error committing transaction. %w", err)
	}

	return nil
}