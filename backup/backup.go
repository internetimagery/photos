package backup

import (
	"log"
	"path/filepath"

	"github.com/internetimagery/photos/context"
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
