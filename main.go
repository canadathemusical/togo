package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/jessevdk/go-flags"
	_ "github.com/mattn/go-sqlite3"
)

type Options struct {
	List ListCommand `command:"list" description:"List tasks"`
	Add  AddCommand  `command:"add" description:"Add a new task"`
	// Edit      EditCommand      `command:"edit" description:"Edit an existing task"`
	// Delete    DeleteCommand    `command:"delete" description:"Delete a task"`
	// CloudSave CloudSaveCommand `command:"cloudsave" description:"Save tasks to the cloud"`
	Done DoneCommand `command:"done" description:"Mark a task as done"`
}

type DoneCommand struct {
	ID string `positional-arg-name:"id"`
}

type AddCommand struct {
	Title string `positional-arg-name:"title"`
}

type ListCommand struct{}

func (cmd *ListCommand) Execute(args []string) error {
	db, err := sql.Open("sqlite3", "./todo.db")
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

func (cmd *AddCommand) Execute(args []string) error {
	// how do I get the argument after add?

	status := "ACTIVE"
	mode := "default"
	notes := ""

	fmt.Println(cmd.Title)
	db, err := sql.Open("sqlite3", "./todo.db")
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

func (cmd *DoneCommand) Execute(args []string) error {
	// how do I get the argument after done?

	db, err := sql.Open("sqlite3", "./todo.db")
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

func main() {
	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        notes TEXT,
        status TEXT,
        mode TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}

	// if there are no args, run list command
	//	if len(os.Args) == 1 {
	//		for _, arg := range os.Args {
	//			fmt.Println(arg)
	//		}
	//	}

	// Set up command parsing
	parser := flags.NewParser(&Options{}, flags.Default)
	_, err = parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}

func wrapText(text string) (string, error) {
	// Get the terminal size
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return "", fmt.Errorf("error getting terminal size: %w", err)
	}

	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)

	lineLength := 0
	for scanner.Scan() {
		word := scanner.Text()
		if lineLength+len(word)+1 > width {
			result.WriteString("\n")
			lineLength = 0
		}
		if lineLength > 0 {
			result.WriteString(" ")
			lineLength++
		}
		result.WriteString(word)
		lineLength += len(word)
	}

	return result.String(), nil
}
