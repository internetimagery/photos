// Add command
package main

import (
	"flag"
	"fmt"
)

var test bool

func cmd_add_init() *command {
	set := flag.NewFlagSet("add", flag.ExitOnError)
	set.BoolVar(&test, "test", false, "Testing command")
	return &command{Set: set, Run: cmd_add_run}
}

func cmd_add_run(args []string) error {
	fmt.Println("RUNNING ADD!", test)
	return nil
}
