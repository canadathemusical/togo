package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Person struct to hold data
type Person struct {
	ID   int
	Name string
	Age  int
}

// Command-line options and subcommands
type Options struct {
	Create CreateCommand `command:"create" description:"Create a new person"`
	List   ListCommand   `command:"list" description:"List all people"`
	Get    GetCommand    `command:"get" description:"Get a person by ID"`
	Update UpdateCommand `command:"update" description:"Update a person by ID"`
	Delete DeleteCommand `command:"delete" description:"Delete a person by ID"`
}

// CreateCommand for creating a new person
type CreateCommand struct {
	Name string `short:"n" long:"name" description:"Name of the person" required:"true"`
	Age  int    `short:"a" long:"age" description:"Age of the person" required:"true"`
}

// ListCommand for listing all people
type ListCommand struct{}

// GetCommand for retrieving a person by ID
type GetCommand struct {
	ID int `short:"i" long:"id" description:"ID of the person" required:"true"`
}

// UpdateCommand for updating a person's details
type UpdateCommand struct {
	ID   int    `short:"i" long:"id" description:"ID of the person" required:"true"`
	Name string `short:"n" long:"name" description:"New name of the person" required:"true"`
	Age  int    `short:"a" long:"age" description:"New age of the person" required:"true"`
}

// DeleteCommand for deleting a person by ID
type DeleteCommand struct {
	ID int `short:"i" long:"id" description:"ID of the person" required:"true"`
}

func main() {
	// Open or create SQLite database
	var err error
	db, err = sql.Open("sqlite3", "./people.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	createTable()

	// Set up command parsing
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)

	// Parse the arguments and execute the corresponding command
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}

// Create a table
func createTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS people (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"age" INTEGER
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

// Execute method for CreateCommand
func (c *CreateCommand) Execute(args []string) error {
	stmt, err := db.Prepare("INSERT INTO people(name, age) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(c.Name, c.Age)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New person created successfully")
	return nil
}

// Execute method for ListCommand
func (l *ListCommand) Execute(args []string) error {
	rows, err := db.Query("SELECT id, name, age FROM people")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("List of people:")
	for rows.Next() {
		var person Person
		err := rows.Scan(&person.ID, &person.Name, &person.Age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", person.ID, person.Name, person.Age)
	}
	return nil
}

// Execute method for GetCommand
func (g *GetCommand) Execute(args []string) error {
	row := db.QueryRow("SELECT id, name, age FROM people WHERE id = ?", g.ID)

	var person Person
	err := row.Scan(&person.ID, &person.Name, &person.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Person not found")
			return nil
		}
		log.Fatal(err)
	}

	fmt.Printf("ID: %d, Name: %s, Age: %d\n", person.ID, person.Name, person.Age)
	return nil
}

// Execute method for UpdateCommand
func (u *UpdateCommand) Execute(args []string) error {
	stmt, err := db.Prepare("UPDATE people SET name = ?, age = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(u.Name, u.Age, u.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Person with ID %d updated successfully\n", u.ID)
	return nil
}

// Execute method for DeleteCommand
func (d *DeleteCommand) Execute(args []string) error {
	stmt, err := db.Prepare("DELETE FROM people WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(d.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Person with ID %d deleted successfully\n", d.ID)
	return nil
}
