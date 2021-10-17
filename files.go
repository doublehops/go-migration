package migrations

import (
	"fmt"
	"os"
	"path/filepath"
)

// listMigrationFiles will get migration files from the configured path.
func (h *Handle) listMigrationFiles() ([]string, error) {

	var files []string
	err := filepath.Walk(h.path, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return files, fmt.Errorf("unable to find migration files. %w", err)
	}

	return files, nil
}
