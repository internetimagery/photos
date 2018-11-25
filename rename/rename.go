package rename

import (
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/copy"
	"github.com/internetimagery/photos/sort"

	"github.com/internetimagery/photos/format"
	"github.com/internetimagery/photos/lock"
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

	// Get event from path
	event := format.NewEvent(cxt.WorkingDir)

	// Get source path
	sourcePath := filepath.Join(cxt.WorkingDir, SOURCEDIR)

	// Grab files from given path
	mediaList, err := event.GetMedia()
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
			media.Event = event.Name
			date, err := sort.GetMediaDate(media.Path)
			if err != nil {
				return err
			}
			media.Date = &date
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

	//////////// Now make some changes! /////////////

	// Make source file directory if it doesn't exist
	if err = os.Mkdir(sourcePath, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	// Run through files!
	// TODO: Make this happen in parallel
	for src, dest := range renameMap {
		tempDest := format.MakeTempPath(src) // Temporary file to create before calling it complete.

		log.Println("Renaming:", src)

		// Create a placeholder file to lock in the spot
		handle, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return err
		}
		handle.Close()
		defer func() { // Cleanup
			if err != nil {
				os.Remove(dest)
			}
		}()

		// Create environment for command
		setEnvironment(src, tempDest, cxt)

		if compress {

			// Grab compress command or use a default command. Do the compression.
			command := cxt.Config.Compress.GetCommand(src)
			if command == "" {
				// We have no command. Just copy the file across directly
				log.Println("Copying:", src)
				if err = <-copy.File(src, tempDest); err != nil {
					return err
				}
			} else {
				// We have a command. Prep and execute it.
				log.Println("Compressing:", src)
				var com *exec.Cmd
				com, err = cxt.PrepCommand(command)
				if err != nil {
					return err
				}
				log.Println("Running:", com.Args)
				if err = com.Run(); err != nil {
					return err
				}

				// Verify file made it to its location and it matches
				var desthandle, srchandle *os.File
				desthandle, err = os.Open(tempDest)
				if err != nil {
					return err
				}
				var srchash, desthash string
				desthash, err = lock.GeneratePerceptualHash("difference", desthandle)
				desthandle.Close()
				if err == nil {
					srchandle, err = os.Open(src)
					if err != nil {
						return err
					}
					srchash, err = lock.GeneratePerceptualHash("difference", srchandle)
					srchandle.Close()
					if err == nil {
						var issame bool
						if issame, err = lock.IsSamePerceptualHash(desthash, srchash); err == nil && !issame {
							err = fmt.Errorf("Compressed image does not match source '%s", src)
							return err
						}
					}
				}
				if err != nil && err != image.ErrFormat {
					return err
				}
			}
		} else {
			// We asked not to compress the file. Just copy it instead
			log.Println("Copying:", src)
			if err = <-copy.File(src, tempDest); err != nil {
				return err
			}
		}

		// Move file to its correct location
		if err = os.Rename(tempDest, dest); err != nil {
			return err
		}

		// Move source file to source folder.
		if err = os.Rename(src, sort.UniqueName(sourceMap[src])); err != nil {
			return err
		}
	}
	return nil
}
