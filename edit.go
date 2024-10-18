package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"

	_ "github.com/jessevdk/go-flags"
	_ "github.com/mattn/go-sqlite3"
)

type EditCommand struct {
	ID string `positional-arg-name:"id"`
}

func (cmd *EditCommand) Execute(args []string) error {
	dbPath, err := getDBPath()
	if err != nil {
		panic(err)
	}

	if len(args) == 0 {
		fmt.Println("No task ID provided")
		return nil
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sqlStmt := `UPDATE tasks SET title = ? WHERE id = ?`
	for _, id := range args {
		title, err := getTaskByID(args[0])
		if err != nil {
			panic(err)
		}

		newTitle, err := editTask(title)
		if err != nil {
			panic(err)
		}
		if title == newTitle {
			fmt.Println("No changes made")
			return nil
		}

		_, err = db.Exec(sqlStmt, newTitle, id)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func editTask(title string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	tmpfile, err := os.CreateTemp("", "task")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name())
	_, err = tmpfile.WriteString(title)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return "", err
	}
	return string(content), nil
}
