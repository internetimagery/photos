// // Main entry point.
//
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/internetimagery/photos/backup"
	"github.com/internetimagery/photos/config"
	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/rename"
	"github.com/internetimagery/photos/sort"
	"github.com/internetimagery/photos/tags"
)

// VERSION : Version information
const VERSION = "0.3"

// sendHelp : Print out helpful message.
func sendHelp() {
	fmt.Println("Command to manage photos naming, compression, backup.")
	fmt.Println("Usage:")
	root := filepath.Base(os.Args[0])
	fmt.Println("  ", root, "version                                   ", "// Print out current version of the tool.")
	fmt.Println("  ", root, "init <name>                               ", "// Set up a new project. Creates a config file also serving as the root of the project.")
	fmt.Println("  ", root, "sort <filename> <filename> ...            ", "// Bring in external files, and sort them by date.")
	fmt.Println("  ", root, "rename                                    ", "// Rename (and compress) files in current directory to their parent directory's namespace (event).")
	fmt.Println("  ", root, "tag <filename> [--remove] <tag> <tag> ... ", "// Add and optionally remove tags from renamed files.")
	fmt.Println("  ", root, "backup <name>                             ", "// Execute specified procedure in config to backup files from the current directory.")
}

// question : Ask yes or no
func question() bool {
	fmt.Print("Is this ok? (y|n) : ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil && err != io.EOF {
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
				return fmt.Errorf("Please provide a name for your project.")
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
			return fmt.Errorf("Already within a project. Cannot initialize...")
		} else {
			return err
		}
		return nil
	}

	// Handle being outside project. Common error across the rest of the functions
	if os.IsNotExist(err) {
		return fmt.Errorf("Project has not been set up. Run 'init' to do an intial setup, then add commands to the file created.")
	} else if err != nil {
		return err
	}

	// Nab the rest of the commands
	switch args[1] {

	case "sort": // Sort files in the working directory into folders of their date
		if len(args) < 3 {
			return fmt.Errorf("Please provide a source directory to sort")
		}
		fmt.Printf("About to sort media in '%s'\n", cxt.WorkingDir)
		if question() {
			fmt.Println("Sorting...")
			if err = sort.SortMedia(cxt, args[2:]...); err != nil {
				return err
			}
		}

	case "rename": // Rename files (and optionally compress them) within working directory
		if cxt.WorkingDir == cxt.Root {
			return fmt.Errorf("Cannot rename media in the root directory (same place as config file.)")
		}
		fmt.Printf("About to rename media in '%s'\n", cxt.WorkingDir)
		if question() {
			fmt.Printf("Renaming media in '%s'\n", cxt.WorkingDir)
			// TODO: Add --no-compress option
			if err = rename.Rename(cxt, true); err != nil {
				return err
			}
		}

	case "tag": // Tag files. Assist searching etc.
		if len(args) < 4 {
			return fmt.Errorf("Please provide a filename, and some tags")
		}
		// Validate and collect options
		remove := false
		filename := args[2]
		if !filepath.IsAbs(filename) { // Could use filepath.Abs, but want to be able to test
			filename = filepath.Join(cxt.WorkingDir, filename)
			if !filepath.HasPrefix(filename, cxt.Root) {
				return fmt.Errorf("Path is outside project '%s'", filename)
			}
		}

		tagNames := []string{}
		for _, arg := range args[3:] {
			if strings.HasPrefix(arg, "-") {
				if arg == "--remove" {
					remove = true
				} else {
					return fmt.Errorf("Unrecognised option '%s'. Did you mean --rename?", arg)
				}
			} else {
				tagNames = append(tagNames, arg)
			}
		}

		if remove {
			tags.RemoveTag(filename, tagNames...)
		} else {
			tags.AddTag(filename, tagNames...)
		}

	case "backup": // Backup files within working directory to specified destination
		if len(args) < 3 {
			return fmt.Errorf("please provide a name for the backup script you wish to run")
		}
		fmt.Printf("About to run backup scripts that match the name '%s'.\nTo backup the media in '%s'\n", args[2], cxt.WorkingDir)
		if question() {
			fmt.Printf("Backing up media in '%s'\n", cxt.WorkingDir)
			if err = backup.RunBackup(cxt, args[2]); err != nil {
				return err
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
		fmt.Println(err)
		return
	}
	if err = run(cwd, os.Args); err != nil {
		fmt.Println(err)
	}
}
