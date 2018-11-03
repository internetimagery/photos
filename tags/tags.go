package tags

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/internetimagery/photos/format"
)

// getMedia : Helper to get media, and validate things
func getMedia(filename string) (*format.Media, error) {
	media := format.NewMedia(filename)
	if info, err := os.Stat(filename); err != nil {
		return media, err
	} else if !info.Mode().IsRegular() {
		return media, fmt.Errorf("Filepath is not a regular file! '%s'", filename)
	}
	event := filepath.Base(filepath.Dir(filename))
	if media.Event != event { // Media isn't in the event, set index to 0
		media.Index = 0
	}
	return media, nil
}

// AddTag : Apply tagnames to a file
func AddTag(filename string, tagnames ...string) error {
	media, err := getMedia(filename)
	if err != nil {
		return err
	}
	if media.Index == 0 { // Media not formatted. Leave it alone
		return nil
	}
	oldname, err := media.FormatName()
	if err != nil {
		return err
	}
	// Apply tags
	for _, tagname := range tagnames {
		tagname = strings.TrimSpace(tagname)
		if tagname != "" {
			media.Tags[tagname] = struct{}{}
		}
	}
	// Build new path
	newname, err := media.FormatName()
	if err != nil {
		return err
	} else if oldname == newname { // Nothing has changed. Nothing to do...
		return nil
	}
	fileDir := filepath.Dir(filename)
	newPath := filepath.Join(fileDir, newname)
	// Ensure newpath does not exist
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		if err == nil {
			return os.ErrExist
		}
		return err
	}
	return os.Rename(filename, newPath)
}

// RemoveTag : Remove tagnames from a file
func RemoveTag(filename string, tagnames ...string) error {
	media, err := getMedia(filename)
	if err != nil {
		return err
	}
	if media.Index == 0 { // Media not formatted. Leave it alone
		return nil
	}
	oldname, err := media.FormatName()
	if err != nil {
		return err
	}
	// Remove tags
	for _, tagname := range tagnames {
		delete(media.Tags, tagname)
	}
	newname, err := media.FormatName()
	if err != nil {
		return err
	} else if oldname == newname { // No change. We're done here!
		return nil
	}
	// Build new path
	fileDir := filepath.Dir(filename)
	newPath := filepath.Join(fileDir, newname)
	// Ensure newpath does not exist
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		if err == nil {
			return os.ErrExist
		}
		return err
	}
	return os.Rename(filename, newPath)
}
