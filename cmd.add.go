// Add command
package main

import (
	"fmt"
)

var test bool

func cmd_add_init() *Command {
	com := NewCommand("add", cmd_add_run)
	com.Set.BoolVar(&test, "test", false, "Testing command")
	return com
}

func cmd_add_run(args []string) error {
	fmt.Println("RUNNING ADD!", test)
	return nil
}
