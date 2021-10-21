package action

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// CreateMigration will copy template file into new fil
func (a *Action) CreateMigration(path string) error {
	currentTime := time.Now()
	curTime := fmt.Sprintf(currentTime.Format("20060102_150405_"))
	name := curTime + a.Name + ".json"
	path = path + "/" + name

	template := `{
  "up": [
    "CREATE TABLE test ( name VARCHAR(255))"
  ],
  "down": [
    "DROP TABLE test"
  ]
}`

	err := ioutil.WriteFile(path, []byte(template), 0644)
	if err != nil {
		return fmt.Errorf("unable to write template file. %s", err)
	}

	os.Stderr.WriteString("Migration file created: " + name + "\n")

	return nil
}
