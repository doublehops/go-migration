package action

import "database/sql"

type Action struct {
	Action string
	Number int
	Name   string

	DB *sql.DB
}

type File struct {
	Filename string
	Queries  *Queries
}

type Queries struct {
	Up   []string `json:"up"`
	Down []string `json:"down"`
}