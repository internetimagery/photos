// // Main entry point.
//
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/internetimagery/photos/backup"
	"github.com/internetimagery/photos/config"
	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/rename"
	"github.com/internetimagery/photos/sort"
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

// run : Do the thing
func run(cwd string, args []string) error {
	// Check for no arguments
	if len(args) == 1 {
		sendHelp()
		return nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	cxt, err := context.NewContext(cwd)

	// We have an argument, nab it and do stuff!

	// Start with special cases
	switch args[1] {

	case "-h":
		sendHelp()
		return nil

	case "--help":
		sendHelp()
		return nil

	case "-v":
		fmt.Println(VERSION)
		return nil

	case "version":
		fmt.Println(VERSION)
		return nil

	case "init": // Create a starter config file at working directory, to signify the root of the project.
		if os.IsNotExist(err) {
			if len(args) < 3 {
				fmt.Println("Please provide a name for your project.")
			} else {
				name := args[2]
				fmt.Printf("About to initialize your project '%s' in '%s'\n", name, cwd)
				if question() {
					configPath := filepath.Join(cwd, context.ROOTCONF)
					newConfig := config.NewConfig(name)
					fmt.Printf("Creating config file '%s'\n", configPath)
					fmt.Println("Be sure to edit it later with what you need. :)")
					handle, err := os.Create(configPath)
					if err != nil {
						return err
					}
					defer handle.Close()
					newConfig.Save(handle)
				}
			}
		} else if err == nil {
			fmt.Println("Already within a project. Cannot initialize...")
		} else {
			return err
		}
		return nil
	}

	// Handle being outside project. Common error across the rest of the functions
	if os.IsNotExist(err) {
		fmt.Println("Project has not been set up. Run 'init' to do an intial setup, then add commands to the file created.")
		return nil
	} else if err != nil {
		return err
	}

	// Nab the rest of the commands
	switch args[1] {

	case "sort": // Sort files in the working directory into folders of their date
		if cxt.WorkingDir == cxt.Root {
			fmt.Println("Cannot run Sort in the root directory (same place as config file.)")
		} else {
			fmt.Printf("About to sort media in '%s'\n", cxt.WorkingDir)
			if question() {
				fmt.Println("Sorting...")
				if err = sort.SortMedia(cxt); err != nil {
					return err
				}
			}
		}

	case "rename": // Rename files (and optionally compress them) within working directory
		if cxt.WorkingDir == cxt.Root {
			fmt.Println("Cannot rename media in the root directory (same place as config file.)")
		} else {
			fmt.Printf("About to rename media in '%s'\n", cxt.WorkingDir)
			if question() {
				fmt.Printf("Renaming media in '%s'\n", cxt.WorkingDir)
				// TODO: Add --no-compress option
				if err = rename.Rename(cxt, true); err != nil {
					return err
				}
			}
		}

	case "backup": // Backup files within working directory to specified destination
		if len(args) < 3 {
			fmt.Println("Please provide a name for the backup script you wish to run.")
		} else {
			fmt.Printf("About to run backup scripts that match the name '%s'.\nTo backup the media in '%s'\n", args[2], cxt.WorkingDir)
			if question() {
				fmt.Printf("Backing up media in '%s'\n", cxt.WorkingDir)
				if err = backup.RunBackup(cxt, args[2]); err != nil {
					return err
				}
			}
		}

	default:
		fmt.Println("Unrecognized command", args[1])
		sendHelp()
	}
	return nil
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err = run(cwd, os.Args); err != nil {
		panic(err)
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
