package migrations

import (
	"database/sql"
	"fmt"
)

type MigrationRecord struct {
	Filename string
	CreatedAt string
}

// getLatestRanMigration will find the last processed migration.
func (h *Handle) getLatestRanMigration() (string, error) {

	var record MigrationRecord
	rows, err := h.db.Query(GetLatestMigrationSQL)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(">> NO ROWS FOUND")
			return "", nil
		}
	}
	defer rows.Close()
	if rows == nil {
		return "", nil
	}

	for rows.Next() {
		if err = rows.Scan(&record.Filename, &record.CreatedAt); err != nil {
			fmt.Printf(">> Query: %s\n", GetLatestMigrationSQL)
			return "", fmt.Errorf("error retrieving rows from migration table. %s", err)
		}
	}

	return record.Filename, nil
}

// getMigrationsNotRun will find the last processed migration.
func (h *Handle) getMigrationsNotRun(files []string, lastRanMigration string) []string {
	var migrationsNotRan []string
	var foundLastRan bool = false

	for _, file := range files {
		if file == lastRanMigration {
			foundLastRan = true
			continue
		}
		if !foundLastRan {
			continue
		}

		migrationsNotRan = append(migrationsNotRan, file)
	}

	return migrationsNotRan
}

// addMigrationTable will add a `migration` table to the database to track what has been
func (h *Handle) addMigrationTable() error {

	_, err := h.db.Exec(CreateMigrationsTable)
	if err != nil {
		return fmt.Errorf("error creating migrations table. %s", err)
	}

	return nil
}
