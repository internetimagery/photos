package rename

import (
	"fmt"

	"github.com/internetimagery/photos/context"

	"github.com/internetimagery/photos/format"
)

// Rename : Rename and compress files within an event (directory)
func Rename(directoryPath string, cxt *context.Context) error {
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

	fmt.Println("Max Index", maxIndex)
	return nil
}
