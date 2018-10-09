package rename

import (
	"fmt"
	"path/filepath"

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

	fmt.Println(renameMap)
	fmt.Println(sourceMap)
	return nil
}
