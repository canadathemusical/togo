package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/term"

	"github.com/jessevdk/go-flags"
	_ "github.com/mattn/go-sqlite3"
)

type Options struct {
	List   ListCommand   `command:"list" description:"List tasks"`
	Add    AddCommand    `command:"add" description:"Add a new task"`
	Edit   EditCommand   `command:"edit" description:"Edit an existing task"`
	Delete DeleteCommand `command:"delete" description:"Delete a task"`
	// CloudSave CloudSaveCommand `command:"cloudsave" description:"Save tasks to the cloud"`
	Done DoneCommand `command:"done" description:"Mark a task as done"`
}

func main() {
	dbPath, err := getDBPath()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", dbPath)
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

	// if no flags have been used, show the list of tasks
	if len(os.Args) == 1 {
		cmd := &ListCommand{}
		err := cmd.Execute(nil)
		if err != nil {
			panic(err)
		}
		return
	}

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

// Function to get the configuration directory based on OS
func getConfigDir() (string, error) {
	var configDir string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		configDir = filepath.Join(homeDir, "AppData", "Local", "togo")
	case "darwin":
		configDir = filepath.Join(homeDir, "Library", "Application Support", "togo")
	default: // Unix-like systems
		configDir = filepath.Join(homeDir, ".config", "togo")
	}

	return configDir, nil
}

// Function to get the path to the database file
func getDBPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	// Ensure the directory exists
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("error creating config directory: %w", err)
	}

	return filepath.Join(configDir, "todo.db"), nil
}
