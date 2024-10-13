package main

import (
	// "flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/jessevdk/go-flags"
)

// func printList() {
// 	fmt.Println("list")
// }

type addValue struct {
	Value string
	Bool  bool
}

func (v *addValue) Set(s string) error {
	v.Value = s
	v.Bool = true
	return nil
}

func (v *addValue) String() string {
	return v.Value
}

func (v *addValue) IsBoolFlag() bool {
	return true
}

func main() {
	var opts struct {
		Add  addValue `long:"add" short:"a" description:"add a todo item"`
		List bool     `long:"list" short:"l" description:"list all todo items"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		// Handle error
	}

	if opts.List {
		fmt.Println("list")
	}

	if opts.Add.Bool {
		fmt.Println("add")
		if opts.Add.Value == "" {
			editorAdd()
		} else {
			fmt.Println(opts.Add.Value)
		}
	}
}

func editorAdd() {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "example.*.txt") //ioutil.TempFile("", "example.*.txt")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	// Open the editor (vim or nano)
	cmd := exec.Command("vim", tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// Read the content of the temporary file
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		panic(err)
	}

	println(string(content))
}
