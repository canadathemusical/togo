package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DoneCommand struct {
	ID string `positional-arg-name:"id"`
}

func (cmd *DoneCommand) Execute(args []string) error {
	dbPath, err := getDBPath()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, id := range args {
		stmt, err := db.Prepare("UPDATE tasks SET status = ? WHERE id = ?")
		if err != nil {
			panic(err)
		}
		_, err = stmt.Exec("DONE", id)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
