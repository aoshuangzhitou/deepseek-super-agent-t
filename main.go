// Command hello prints a friendly greeting for a human.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// defaultName is used when the caller does not supply a name.
const defaultName = "Human"

// greeting builds the greeting line for name. An empty name falls back to
// defaultName so the message always addresses someone.
func greeting(name string) string {
	if name == "" {
		name = defaultName
	}
	return fmt.Sprintf("Hello, %s!", name)
}

// run writes the greeting to w and reports whether a usage error occurred.
func run(w io.Writer, args []string) error {
	fs := flag.NewFlagSet("hello", flag.ContinueOnError)
	fs.SetOutput(w)
	name := fs.String("name", "", "name to greet (default \""+defaultName+"\")")
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, err := fmt.Fprintln(w, greeting(*name))
	return err
}

func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
