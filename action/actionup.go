package action

import (
	"fmt"
	"github.com/doublehops/go-migration/helpers"
)

// MigrateUp will run new migration/s.
func (a *Action) MigrateUp(migrationFiles []File) error {
	var err error

	for _, file := range migrationFiles {
		err = a.processFileUp(file)
		if err != nil {
			return err
		}
	}

	return nil
}

// processFile will process the queries in the given file. It will attempt to rollback when there is
// and error in one of the queries.
func (a *Action) processFileUp(file File) error {

	tx, err := a.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction. %w", err)
	}
	defer tx.Rollback() // nolint

	helpers.PrintMsg(fmt.Sprintf("Migrating queries from: %s", file.Filename))
	for _, query := range file.Queries {
		_, err = tx.Exec(query)
		if err != nil {
			return fmt.Errorf("\nthere was an error executing query. File: %s; query; %s; Error: %s", file.Filename, file.Queries, err)
		}
	}

	_, err = tx.Exec(InsertMigrationRecordIntoTableSQL, file.Filename)
	if err != nil {
		return fmt.Errorf("unable to update migration table with newly ran migration record. %w\n", err)
	}
	helpers.PrintMsg(" - Success\n")

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("there was an error committing transaction. %w", err)
	}

	return nil
}
