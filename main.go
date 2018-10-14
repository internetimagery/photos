// // Main entry point.
//
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/internetimagery/photos/config"
	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/rename"
)

// VERSION : Version information
const VERSION = "0.0.1"

// sendHelp : Print out helpful message.
func sendHelp() {
	fmt.Println("Command to manage photos naming, compression, backup.")
	fmt.Println("Usage:")
	root := filepath.Base(os.Args[0])
	fmt.Println("    ", root, "init <name>  ", "// Set up a new project. Creates a config file also serving as the root of the project.")
	fmt.Println("    ", root, "sort         ", "// Sort files in the current directory, into folders named after their dates.")
	fmt.Println("    ", root, "rename       ", "// Rename (and compress) files in current directory to their parent directory's namespace (event).")
	fmt.Println("    ", root, "backup <name>", "// Execute specified procedure in config to backup files from the current directory.")
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

	case "-v":
		fmt.Println(VERSION)
		return

	case "version":
		fmt.Println(VERSION)
		return

	case "init": // Create a starter config file at working directory, to signify the root of the project.
		if os.IsNotExist(err) {
			if len(os.Args) < 3 {
				fmt.Println("Please provide a name for your project.")
			} else {
				name := os.Args[2]
				fmt.Printf("About to initialize your project '%s' in '%s'\n", name, cwd)
				if question() {
					configPath := filepath.Join(cwd, context.ROOTCONF)
					newConfig := config.NewConfig(name)
					fmt.Printf("Creating config file '%s'\n", configPath)
					fmt.Println("Be sure to edit it later with what you need. :)")
					handle, err := os.Create(configPath)
					if err != nil {
						panic(err)
					}
					defer handle.Close()
					newConfig.Save(handle)
				}
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

	case "sort": // Sort files in the working directory into folders of their date
		fmt.Println("Sorting files sometime")

	case "rename": // Rename files (and optionally compress them) within working directory
		if cxt.WorkingDir == cxt.Root {
			fmt.Println("Cannot rename media in the root directory (same place as config file.)")
		} else {
			fmt.Printf("About to rename media in '%s'\n", cxt.WorkingDir)
			if question() {
				fmt.Printf("Renaming media in '%s'\n", cxt.WorkingDir)
				if err = rename.Rename(cxt, true); err != nil {
					panic(err)
				}
			}
		}

	case "backup": // Backup files within working directory to specified destination
		fmt.Println("When this is functioning, life will be grand!")

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
