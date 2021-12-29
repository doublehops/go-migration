package go_migration

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/doublehops/go-migration/helpers"
)

// CreateMigration will copy template file into new fil
func (a *Action) CreateMigration(path string) error {
	currentTime := time.Now()
	curTime := currentTime.Format("20060102_150405_")
	name := curTime + a.Name
	upName := name + ".up.sql"
	downName := name + ".down.sql"
	upPath := path + "/" + upName
	downPath := path + "/" + downName

	separatorMessage := fmt.Sprintf("-- You need to separate multiple queries with this dotted line: %s\n\n", QuerySeparator)

	exampleUp :=
`CREATE TABLE news (
    id INT(11) NOT NULL,
    currency_id INT(11) NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    PRIMARY KEY (id),
    FOREIGN KEY (currency_id) REFERENCES currency(id));
	`

	exampleUp = separatorMessage + exampleUp
	exampleDown := separatorMessage +"DROP TABLE news;\n\n"

	err := ioutil.WriteFile(upPath, []byte(exampleUp), 0644)
	if err != nil {
		return fmt.Errorf("unable to write template file: %s. %s", upPath, err)
	}

	err = ioutil.WriteFile(downPath, []byte(exampleDown), 0644)
	if err != nil {
		return fmt.Errorf("unable to write template file: %s. %s", downPath, err)
	}

	helpers.PrintMsg("Migration file created: " + upPath + "\n")
	helpers.PrintMsg("Migration file created: " + downPath + "\n")

	return nil
}
