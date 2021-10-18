package migrations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// listMigrationFiles will get migration files from the configured path.
func (h *Handle) listMigrationFiles() ([]string, error) {

	var files []string
	f, err := ioutil.ReadDir(h.path)
	if err != nil {
		return files, fmt.Errorf("unable to list migration files. %w", err)
	}

	for _, file := range f {
		files = append(files, file.Name())
	}

	return files, nil
}

// getPendingMigrationFiles will loop through all migration files and return the ones that haven't been run yet.
func (h *Handle) getPendingMigrationFiles() ([]string, error) {
	var pendingFiles []string
	var foundLastRan bool = false

	lastRanMigration, err := h.getLatestRanMigration()
	if err != nil {
		return pendingFiles, err
	}
	allFiles, err := h.listMigrationFiles()
	if err != nil {
		return pendingFiles, err
	}

	if lastRanMigration == "" { // No migrations have run yet.
		return allFiles, nil
	}

	for _, file := range allFiles {
		if file == lastRanMigration {
			foundLastRan = true
			continue
		}
		if !foundLastRan {
			continue
		}

		pendingFiles = append(pendingFiles, file)
	}

	return pendingFiles, nil
}

// parseMigrations will iterate through the files and unmarshal the JSON and add to the files slice.
func (h *Handle) parseMigrations(filesToParse []string) ([]File, error) {
	var files []File
	for _, file := range filesToParse {

		thisFile := File{Filename: file}
		data, err := os.ReadFile(h.path+"/"+file)
		if err != nil {
			return files, fmt.Errorf("unable to read file: %s. %s", file, err)
		}

		var q Queries
		err = json.Unmarshal(data, &q)
		if err != nil {
			return files, fmt.Errorf("unable to unmarshal query. %s", err)
		}

		thisFile.Queries = &q
		files = append(files, thisFile)
	}

	return files, nil
}