package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jessevdk/go-flags"
	_ "github.com/mattn/go-sqlite3"
)

type ListCommand struct{}

func (cmd *ListCommand) Execute(args []string) error {
	dbPath, err := getDBPath()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStmt := `SELECT id, title, status FROM tasks WHERE status = "ACTIVE"`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	//  if there are no active tasks print "you have nothing to do"
	activeTasks := false
	for rows.Next() {
		var id int
		var title, status string
		err := rows.Scan(&id, &title, &status)
		if err != nil {
			panic(err)
		}
		if status == "ACTIVE" {
			text, err := wrapText(fmt.Sprintf("[%d] %s \n", id, title))
			if err != nil {
				panic(err)
			}
			fmt.Println(text)
		}
		activeTasks = true
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	if !activeTasks {
		text, err := wrapText("You have nothing to do")
		if err != nil {
			panic(err)
		}
		fmt.Println(text)

	}
	return nil
}
