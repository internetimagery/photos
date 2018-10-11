package rename

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	shlex "github.com/google/shlex"
	"github.com/internetimagery/photos/context"

	"github.com/internetimagery/photos/format"
)

// SOURCEDIR : File to store originals for manual checking
const SOURCEDIR = "Source Media Please check before removing"

// Rename : Rename and compress files within an event (directory)
func Rename(directoryPath string, cxt *context.Context) error {
	// Get event name from path
	eventName := filepath.Base(directoryPath)

	// Get source path
	sourcePath := filepath.Join(directoryPath, SOURCEDIR)

	// Grab files from given path
	mediaList, err := format.GetMediaFromDirectory(directoryPath)
	if err != nil {
		return err
	}

	// Get max index
	maxIndex := 0
	for _, media := range mediaList {
		if maxIndex < media.Index {
			maxIndex = media.Index
		}
	}

	// Map old names to new names
	renameMap := make(map[string]string)
	// Map renames to source
	sourceMap := make(map[string]string)
	for _, media := range mediaList {
		if media.Index == 0 { // Media is not already named correctly
			maxIndex++
			media.Index = maxIndex
			media.Event = eventName
			newName, err := media.FormatName()
			if err != nil {
				return err
			}
			renameMap[media.Path] = filepath.Join(directoryPath, newName)
			sourceMap[media.Path] = filepath.Join(sourcePath, filepath.Base(media.Path))
		}
	}

	// Make sure we actually have something to do
	if len(renameMap) == 0 {
		return nil
	}

	// Check files aren't already in the source directory
	for _, source := range sourceMap {
		if _, err = os.Stat(source); !os.IsNotExist(err) {
			return fmt.Errorf("File already exists: '%s'", source)
		}
	}

	//////////// Now make some changes! /////////////

	// Make source file directory if it doesn't exist
	if err = os.Mkdir(sourcePath, 755); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Run through files!
	for src, dest := range renameMap {

		// Grab compress command or use a default command. Do the compression.
		command := cxt.Config.Compress.GetCommand(src)
		if command == "" {
			command = `cp -a "$SOURCEPATH" "$DESTPATH"`
		}
		env := map[string]string{
			"SOURCEPATH":  src,            // From where are we going?
			"DESTPATH":    dest,           // To where are we headed?
			"ROOTPATH":    cxt.Root,       // Where is the root of our project?
			"WORKINGPATH": cxt.WorkingDir, // Where are we working right now?
		}
		if err = runCommand(command, env); err != nil {
			return err
		}

		// Verify file made it to its location
		if _, err = os.Stat(dest); err != nil {
			return err
		}

		// Move source file to source folder.
		if err = os.Rename(src, sourceMap[src]); err != nil {
			return err
		}
	}
	return nil
}

// runCommand : Helper to run commands, linking outputs to terminal outputs and replacing variables safely
func runCommand(commandString string, environment map[string]string) error {
	buildEnv := func(name string) string {
		return strings.Replace(environment[name], `\`, `\\`, -1)
	}
	commandExpand := os.Expand(commandString, buildEnv)
	commandParts, err := shlex.Split(commandExpand)
	if err != nil {
		return err
	}
	command := exec.Command(commandParts[0], commandParts[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	fmt.Println("Running:", commandParts)
	return command.Run()
}
