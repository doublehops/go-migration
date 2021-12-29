package go_migration

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

// listMigrationFiles will get migration files from the configured path.
func (a *Action) listMigrationFiles() ([]string, error) {

	fileFilter := "."+ a.Action +".sql"

	var files []string
	f, err := ioutil.ReadDir(a.Path)
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
func (a *Action) getPendingMigrationFiles() ([]string, error) {
	var pendingFiles []string
	var foundLastRan = false

	lastRanMigration, err := a.getLatestRanMigration()
	if err != nil {
		return pendingFiles, err
	}
	allFiles, err := a.listMigrationFiles()
	if err != nil {
		return pendingFiles, err
	}

	if lastRanMigration == "" { // No migrations have run yet.
		foundLastRan = true // If no migrations have previously ran, set found as true to start from first file.
	}

	var i = 0
	for _, file := range allFiles {
		if i == a.Number {
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
func (a *Action) getMigrationFilesToRollBack() ([]string, error) {
	var migrationsToRollBack []string
	var foundLastRan = false

	lastRanMigration, err := a.getLatestRanMigration()
	if err != nil {
		return migrationsToRollBack, err
	}
	allFiles, err := a.listMigrationFiles()
	if err != nil {
		return migrationsToRollBack, err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(allFiles)))

	if lastRanMigration == "" { // No migrations have run yet.
		return migrationsToRollBack, nil
	}

	lastRanMigrationShortName := TrimExtension(lastRanMigration)

	var i = 0
	for _, file := range allFiles {
		shortFileName := TrimExtension(file)
		if i == a.Number {
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
func (a *Action) parseMigrations(filesToParse []string) ([]File, error) {
	var files []File
	for _, file := range filesToParse {

		thisFile := File{Filename: file}
		data, err := os.ReadFile(a.Path+"/"+file)
		if err != nil {
			return files, fmt.Errorf("unable to read file: %s. %s", file, err)
		}

		queries := strings.Split(string(data), QuerySeparator)

		thisFile.Queries = queries

		files = append(files, thisFile)
	}

	return files, nil
}