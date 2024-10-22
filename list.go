package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jessevdk/go-flags"
	_ "github.com/mattn/go-sqlite3"
)

type ListCommand struct {
	Args struct {
		Status string `positional-arg-name:"status"`
	} `positional-args:"true"`
}

func (cmd *ListCommand) Execute(args []string) error {
	// dbPath, err := getDBPath()
	// if err != nil {
	// 	panic(err)
	// }

	// db, err := sql.Open("sqlite3", dbPath)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// fmt.Println("status", cmd.Args.Status)

	// sqlStmt := `SELECT id, title, status FROM tasks WHERE status = "ACTIVE"`
	// rows, err := db.Query(sqlStmt)
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()

	// //  if there are no active tasks print "you have nothing to do"
	// activeTasks := false
	// for rows.Next() {
	// 	var id int
	// 	var title, status string
	// 	err := rows.Scan(&id, &title, &status)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if status == "ACTIVE" {
	// 		text, err := wrapText(fmt.Sprintf("[%d] %s \n", id, title))
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		fmt.Println(text)
	// 	}
	// 	activeTasks = true
	// }
	// err = rows.Err()
	// if err != nil {
	// 	panic(err)
	// }
	// if !activeTasks {
	// 	text, err := wrapText("You have nothing to do")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(text)

	// }
	// return nil
	subcommand := strings.ToUpper(cmd.Args.Status)
	switch subcommand {
	case "", "ACTIVE":
		return cmd.listActive()
	case "ALL":
		return cmd.listAll()
	case "DONE":
		return cmd.listDone()
	default:
		return fmt.Errorf("unknown subcommand: %s", cmd.Args.Status)
	}
}

func (cmd *ListCommand) listActive() error {
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

func (cmd *ListCommand) listAll() error {
	fmt.Println("Listing all tasks")
	// Implement listing all tasks
	return nil
}

func (cmd *ListCommand) listDone() error {
	fmt.Println("Listing done tasks")
	// Implement listing done tasks
	return nil
}
