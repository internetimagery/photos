package backup

import (
	"fmt"
	"path/filepath"

	"github.com/internetimagery/photos/context"
)

func RunBackup(cxt *context.Context, name string) error {

	// Prep our environment for command
	relpath, _ := filepath.Rel(cxt.Root, cxt.WorkingDir)
	env := map[string]string{
		"SOURCEPATH": cxt.WorkingDir,
		"ROOTPATH":   cxt.Root,
		"RELPATH":    relpath,
	}

	// Get our command
	commands := cxt.Config.Backup.GetCommands(name, env)
	run := 0
	for _, command := range commands {
		if command != "" {
			run++

			// Run our backup command
			err := context.RunCommand(command)
			if err != nil {
				return err
			}
		}
	}

	// Check if we ran anything at all
	if run == 0 {
		return fmt.Errorf("No commands match the name '%s'", name)
	}
	return nil
}
