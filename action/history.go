package action

import (
	"database/sql"
	"fmt"
)

type MigrationRecord struct {
	ID int
	Filename string
	CreatedAt string
}

// getLatestRanMigration will find the last processed migration.
func (a *Action) getLatestRanMigration() (string, error) {

	var record MigrationRecord
	rows, err := a.DB.Query(GetLatestMigrationSQL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
	}
	defer rows.Close()
	if rows == nil {
		return "", nil
	}

	for rows.Next() {
		if err = rows.Scan(&record.ID, &record.Filename, &record.CreatedAt); err != nil {
			return "", fmt.Errorf("error retrieving rows from migration table. %s", err)
		}
	}

	return record.Filename, nil
}

// addMigrationTable will add a `migration` table to the database to track what has been
func (a *Action) addMigrationTable() error {

	_, err := a.DB.Exec(CreateMigrationsTable)
	if err != nil {
		return fmt.Errorf("error creating migrations table. %s", err)
	}

	return nil
}
