package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"

	_ "github.com/jessevdk/go-flags"
	_ "github.com/mattn/go-sqlite3"
)

type AddCommand struct {
	Title string `positional-arg-name:"title"`
}

func (cmd *AddCommand) Execute(args []string) error {
	dbPath, err := getDBPath()
	if err != nil {
		panic(err)
	}

	// how do I get the argument after add?

	status := "ACTIVE"
	mode := "default"
	notes := ""

	if len(args) == 0 {
		// No arguments, call editorAdd
		title, err := getTaskFromEditor()
		if err != nil {
			panic(err)
		}
		if title == "" {
			fmt.Println("No task added")
			return nil
		}
		args = append(args, title)
	}

	fmt.Println(cmd.Title)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStmt := `INSERT INTO tasks(title, notes, status, mode) VALUES(?, ?, ?, ?)`

	for _, title := range args {

		_, err = db.Exec(sqlStmt, title, notes, status, mode)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func getTaskFromEditor() (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	tmpfile, err := os.CreateTemp("", "task")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name())

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
