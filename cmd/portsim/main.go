// Command portsim is the entry point of the container port simulator.
package main

import (
	"fmt"
	"io"
	"os"
)

const version = "0.0.0-scaffold"

// stubs maps each subcommand that is not built yet to its parent ticket.
var stubs = map[string]string{
	"generate": "#2",
	"run":      "#5",
	"compare":  "#5",
	"serve":    "#6",
	"verify":   "#7",
}

const usage = `usage: portsim <command>

commands: generate, run, compare, serve, verify
flags:    --help, --version
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "--help", "-h", "help":
		fmt.Fprint(stdout, usage)
		return 0
	case "--version":
		fmt.Fprintln(stdout, "portsim "+version)
		return 0
	}
	if ticket, ok := stubs[args[0]]; ok {
		fmt.Fprintf(stderr, "not implemented yet: see %s\n", ticket)
		return 1
	}
	fmt.Fprintf(stderr, "unknown command %q\n%s", args[0], usage)
	return 2
}
