package backup

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/format"
	"github.com/internetimagery/photos/lock"
	"github.com/internetimagery/photos/rename"
)

// setEnvironment : Set up environment variables for the command context
func setEnvironment(cxt *context.Context) {
	relpath, _ := filepath.Rel(cxt.Root, cxt.WorkingDir)
	cxt.Env["SOURCEPATH"] = cxt.WorkingDir
	cxt.Env["ROOTPATH"] = cxt.Root
	cxt.Env["WORKINGPATH"] = cxt.WorkingDir
	cxt.Env["RELPATH"] = filepath.ToSlash(relpath)
}

// RunBackup : Run backup commands given a name. Can accept wildcards to run more than one.
func RunBackup(cxt *context.Context, name string) error {

	sortDir := filepath.Join(cxt.Root, cxt.Config.Sorted)

	// Validate our project
	if err := filepath.Walk(cxt.WorkingDir, func(filename string, info os.FileInfo, err error) error {
		if info.IsDir() { // Lock files in directory! Also a validation
			return lock.LockEvent(filename, false)
		}
		if !format.IsUsable(filename) { // Ignore any file deemed unusable
			return nil
		}
		if strings.Contains(filename, rename.SOURCEDIR) {
			return fmt.Errorf("refusing to backup with source files still inside '%s'", filename)
		}
		if strings.HasPrefix(filename, sortDir) {
			return fmt.Errorf("refusing to backup within the sorting directory '%s'", filename)
		}
		event := filepath.Base(filepath.Dir(filename))
		if media := format.NewMedia(info.Name()); media.Index == 0 || media.Event != event {
			return fmt.Errorf("refusing to backup with unformatted files still inside '%s'", filename)
		}
		return nil
	}); err != nil {
		return err
	}

	// Prep our environment for command
	setEnvironment(cxt)

	// Get our command
	commands := cxt.Config.Backup.GetCommands(name)
	run := 0
	for _, command := range commands {
		if command != "" {
			run++

			// Run our backup command
			com, err := cxt.PrepCommand(command)
			if err != nil {
				return err
			}
			log.Println("Running:", com.Args)
			err = com.Run()
			if err != nil {
				return err
			}
		}
	}

	// Check if we ran anything at all
	if run == 0 {
		log.Printf("No commands match the name '%s'\n", name)
	}
	return nil
}
