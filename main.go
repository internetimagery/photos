// // Main entry point.
//
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/internetimagery/photos/context"
)

// sendHelp : Print out helpful message.
func sendHelp() {
	fmt.Println("Command to manage photos naming, compression, backup.")
	fmt.Println("Usage:")
	root := filepath.Base(os.Args[0])
	fmt.Println("\t", root, "init", "\t\t// Set up a new project. Creates a config file also serving as the root of the project.")
	fmt.Println("\t", root, "rename", "\t\t// Rename (and compress) files in current directory to their parent directory's namespace (event).")
	fmt.Println("\t", root, "backup <name>", "\t// Execute specified procedure in config to backup files from the current directory.")
}

// question : Ask yes or no
func question() bool {
	fmt.Print("Is this ok? (y|n) : ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(response) == "y"
}

func main() {
	// Check for no arguments
	if len(os.Args) == 1 {
		sendHelp()
		return
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cxt, err := context.NewContext(cwd)

	// We have an argument, nab it and do stuff!

	// Start with special cases
	switch os.Args[1] {

	case "-h":
		sendHelp()
		return

	case "--help":
		sendHelp()
		return

	case "init": // Create a starter config file at working directory, to signify the root of the project.
		if os.IsNotExist(err) {
			fmt.Printf("About to initialize your project in '%s'\n", cwd)
			if question() {
				fmt.Println("YAY DO IT")
				// configPath := filepath.Join(cwd, context.ROOTCONF)
			}
		} else if err == nil {
			fmt.Println("Already within a project. Cannot initialize...")
		} else {
			panic(err)
		}
		return
	}

	// Handle being outside project. Common error across the rest of the functions
	if os.IsNotExist(err) {
		fmt.Println("Project has not been set up. Run 'init' to do an intial setup, then add commands to the file created.")
		return
	} else if err != nil {
		panic(err)
	}

	// Nab the rest of the commands
	switch os.Args[1] {

	case "rename": // Rename files (and optionally compress them) within working directory
		fmt.Printf("About to rename media in '%s'", cxt.WorkingDir)
		if question() {
			fmt.Println("Ok gonna do it I guess")
		}

	case "backup": // Backup files within working directory to specified destination
		fmt.Printf("About to rename media in '%s'", cxt.WorkingDir)
		if question() {
			fmt.Println("Ok gonna do it I guess")
		}

	default:
		fmt.Println("Unrecognized command", os.Args[1])
		sendHelp()
	}
}

//
// import (
// 	"flag"
// 	"fmt"
// 	"os"
// )
//
// // Single command
// type Command struct {
// 	Set *flag.FlagSet
// 	Run func([]string) error
// }
//
// func NewCommand(name string, run func([]string) error) *Command {
// 	return &Command{Set: flag.NewFlagSet(name, flag.ExitOnError), Run: run}
// }
//
// func main() {
// 	// Initialize our commands
// 	coms := make(map[string]*Command)
// 	coms["add"] = cmd_add_init()
//
// 	// If no commands are issued. Send help message.
// 	if len(os.Args) < 2 {
// 		fmt.Println("Available commands:")
// 		for c, _ := range coms {
// 			fmt.Println(c)
// 		}
// 		os.Exit(1)
// 	}
//
// 	// Grab requested command
// 	com := coms[os.Args[1]]
// 	if com == nil {
// 		fmt.Println("Command", os.Args[1], "not valid.")
// 		fmt.Println("Valid commands:")
// 		for c, _ := range coms {
// 			fmt.Println(c)
// 		}
// 		os.Exit(1)
// 	}
//
// 	// Parse commands, and run
// 	com.Set.Parse(os.Args[2:])
// 	args := com.Set.Args()
// 	err := com.Run(args)
// 	if err != nil {
// 		panic(err)
// 	}
// }
