package main

import (
	"database/sql"
	"fmt"

	"github.com/manifoldco/promptui"
	_ "github.com/mattn/go-sqlite3"
)

type DeleteCommand struct {
	ID string `positional-arg-name:"id"`
}

func (cmd *DeleteCommand) Execute(args []string) error {
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
		if !confirmDelete(id) {
			fmt.Println("Task not deleted")
			continue
		}
		stmt, err := db.Prepare("DELETE FROM tasks WHERE id = ?")
		if err != nil {
			panic(err)
		}
		_, err = stmt.Exec(id)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func confirmDelete(id string) bool {
	title, err := getTaskByID(id)
	if err != nil {
		fmt.Printf("Error getting task: %v\n", err)
		return false
	}
	fmt.Println("This will permanently delete the task:", title)
	prompt := promptui.Prompt{
		Label:     "Are you sure you want to delete this task?",
		IsConfirm: true,
	}
	_, err = prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return false
	}
	return true
}

func getTaskByID(id string) (string, error) {
	dbPath, err := getDBPath()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStmt := `SELECT title FROM tasks WHERE id = ?`
	row := db.QueryRow(sqlStmt, id)
	var title string
	err = row.Scan(&title)
	if err != nil {
		return "", err
	}
	return title, nil
}
