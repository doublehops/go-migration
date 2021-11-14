package migrations

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/doublehops/go-migration/action"
)

// listMigrationFiles will get migration files from the configured path.
func (h *Handle) listMigrationFiles() ([]string, error) {

	fileFilter := "."+ h.action.Action +".sql"

	var files []string
	f, err := ioutil.ReadDir(h.path)
	if err != nil {
		return files, fmt.Errorf("unable to list migration files. %w", err)
	}

	for _, file := range f {
		if strings.Contains(file.Name(), fileFilter) {
			files = append(files, file.Name())
		}
	}

	return files, nil
}

// getPendingMigrationFiles will loop through all migration files and return the ones that haven't been run yet.
func (h *Handle) getPendingMigrationFiles() ([]string, error) {
	var pendingFiles []string
	var foundLastRan = false

	lastRanMigration, err := h.getLatestRanMigration()
	if err != nil {
		return pendingFiles, err
	}
	allFiles, err := h.listMigrationFiles()
	if err != nil {
		return pendingFiles, err
	}

	if lastRanMigration == "" { // No migrations have run yet.
		foundLastRan = true // If no migrations have previously ran, set found as true to start from first file.
	}

	var i = 0
	for _, file := range allFiles {
		if i == h.action.Number {
			break
		}
		if file == lastRanMigration {
			foundLastRan = true
			continue
		}
		if !foundLastRan {
			continue
		}

		pendingFiles = append(pendingFiles, file)
		i++
	}

	return pendingFiles, nil
}

// getPreviouslyMigratedFiles will loop through all migration files and return the ones that have already been run.
func (h *Handle) getMigrationFilesToRollBack() ([]string, error) {
	var migrationsToRollBack []string
	var foundLastRan = false

	lastRanMigration, err := h.getLatestRanMigration()
	if err != nil {
		return migrationsToRollBack, err
	}
	allFiles, err := h.listMigrationFiles()
	if err != nil {
		return migrationsToRollBack, err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(allFiles)))

	if lastRanMigration == "" { // No migrations have run yet.
		return allFiles, nil
	}

	lastRanMigrationShortName := action.TrimExtension(lastRanMigration)

	var i = 0
	for _, file := range allFiles {
		shortFileName := action.TrimExtension(file)
		if i == h.action.Number {
			break
		}
		if shortFileName == lastRanMigrationShortName {
			foundLastRan = true
			migrationsToRollBack = append(migrationsToRollBack, file)
			i++
			continue
		}
		if !foundLastRan {
			continue
		}

		migrationsToRollBack = append(migrationsToRollBack, file)
		i++
	}

	return migrationsToRollBack, nil
}

// parseMigrations will iterate through the files and unmarshal the JSON and add to the files slice.
func (h *Handle) parseMigrations(filesToParse []string) ([]action.File, error) {
	var files []action.File
	for _, file := range filesToParse {

		thisFile := action.File{Filename: file}
		data, err := os.ReadFile(h.path+"/"+file)
		if err != nil {
			return files, fmt.Errorf("unable to read file: %s. %s", file, err)
		}

		thisFile.Queries = string(data)

		files = append(files, thisFile)
	}

	return files, nil
}