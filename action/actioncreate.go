package action

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
	name := curTime + a.Name + ".json"
	path = path + "/" + name

	template := `{
  "up": [
    "CREATE TABLE news (id INT(11) NOT NULL, currency_id INT(11) NOT NULL, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME, PRIMARY KEY (id), FOREIGN KEY (currency_id) REFERENCES currency(id));"
  ],
  "down": [
    "DROP TABLE news"
  ]
}`

	err := ioutil.WriteFile(path, []byte(template), 0644)
	if err != nil {
		return fmt.Errorf("unable to write template file. %s", err)
	}

	helpers.PrintMsg("Migration file created: " + name + "\n")

	return nil
}
