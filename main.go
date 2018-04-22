// Main entry point.

package main

import (
	"flag"
	"fmt"
	"os"
)

// Single command
type command struct {
	Set *flag.FlagSet
	Run func([]string) error
}

func main() {
	// Initialize our commands
	coms := make(map[string]*command)
	coms["add"] = cmd_add_init()

	// If no commands are issued. Send help message.
	if len(os.Args) < 2 {
		fmt.Println("Available commands:")
		for c, _ := range coms {
			fmt.Println(c)
		}
		os.Exit(1)
	}

	// Grab requested command
	com := coms[os.Args[1]]
	if com == nil {
		fmt.Println("Command", os.Args[1], "not valid.")
		fmt.Println("Valid commands:")
		for c, _ := range coms {
			fmt.Println(c)
		}
		os.Exit(1)
	}

	// Parse commands, and run
	com.Set.Parse(os.Args[2:])
	args := com.Set.Args()
	err := com.Run(args)
	if err != nil {
		panic(err)
	}
}
