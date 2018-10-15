package rename

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/internetimagery/photos/context"

	"github.com/internetimagery/photos/format"
)

// SOURCEDIR : File to store originals for manual checking
const SOURCEDIR = "Source Media - Please check before removing"

// setEnvironment : Set up environment variables for the command context
func setEnvironment(sourcePath, destPath string, cxt *context.Context) {
	cxt.Env["SOURCEPATH"] = sourcePath
	cxt.Env["DESTPATH"] = destPath
	cxt.Env["ROOTPATH"] = cxt.Root
	cxt.Env["WORKINGPATH"] = cxt.WorkingDir
}

// Rename : Rename and compress files within an event (directory). Optionally compress while renaming.
func Rename(cxt *context.Context, compress bool) error {

	// Get event name from path
	eventName := filepath.Base(cxt.WorkingDir)

	// Get source path
	sourcePath := filepath.Join(cxt.WorkingDir, SOURCEDIR)

	// Grab files from given path
	mediaList, err := format.GetMediaFromDirectory(cxt.WorkingDir)
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
			renameMap[media.Path] = filepath.Join(cxt.WorkingDir, newName)
			sourceMap[media.Path] = filepath.Join(sourcePath, filepath.Base(media.Path))
		}
	}

	// Make sure we actually have something to do
	if len(renameMap) == 0 {
		log.Println("Nothing to rename...")
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
	if err = os.Mkdir(sourcePath, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	// Run through files!
	for src, dest := range renameMap {

		log.Println("Renaming:", src)

		// Create environment for command
		setEnvironment(src, dest, cxt)

		if compress {

			// Grab compress command or use a default command. Do the compression.
			command := cxt.Config.Compress.GetCommand(src)
			if command == "" { // We have no command. Just copy the file across
				log.Println("Moving:", src)
				err = os.Link(src, dest)
				if err != nil {
					return err
				}
			} else { // We have a command. Prep and execute it.
				log.Println("Compressing:", src)
				com, err := cxt.PrepCommand(command)
				if err != nil {
					return err
				}
				log.Println("Running:", com.Args)
				if err = com.Run(); err != nil {
					return err
				}
			}
		} else {
			if err = os.Link(src, dest); err != nil {
				return err
			}
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
