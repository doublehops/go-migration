package action

import (
	"database/sql"
	"strings"
)

type Action struct {
	Action string
	Number int
	Name   string

	DB *sql.DB
}

type File struct {
	Filename string
	Queries  string
}

func TrimExtension(filename string) string {
	var str string

	str = strings.Replace(filename, ".up.sql", "", 1)
	str = strings.Replace(str, ".down.sql", "", 1)

	return str
}