package sort

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/format"
	"github.com/internetimagery/photos/copy"
	"github.com/rwcarlsen/goexif/exif"
)

// GetMediaDate : Get modification date, or date taken (EXIF data) from file
func GetMediaDate(filePath string) (time.Time, error) {

	// Get a handle on things... get it!
	handle, err := os.Open(filePath)
	if err != nil {
		return time.Time{}, err
	}
	defer handle.Close()
	info, err := handle.Stat()
	if err != nil {
		return time.Time{}, err
	}

	// Only process regular files
	if !info.Mode().IsRegular() {
		return time.Time{}, fmt.Errorf("can only get media date from files")
	}

	// Try processing exif data
	if exifData, err := exif.Decode(handle); err == nil {
		taken, err := exifData.DateTime()
		if err == nil {
			return taken, nil
		}
		return time.Time{}, err
	}
	// We failed to decode it? Ah well... fall back to using modtime
	return info.ModTime(), nil
}

// FormatDate : Format date into simple YY-MM-DD style
func FormatDate(date time.Time) string {
	return date.Format("06-01-02")
}

// UniqueName : Ensure name is a unique filename so as to not override existing
func UniqueName(filename string) string {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return filename
	}
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	index := 0
	for {
		index++
		filename = fmt.Sprintf("%s_%d%s", name, index, ext)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			break
		}
	}
	return filename
}

// SortMedia : Grab dates assosicated with media in working directory, and place them in corresponding folders
func SortMedia(cxt *context.Context, copyFiles bool, source ...string) error {

	// Validate our inputs
	if len(source) == 0 {
		return fmt.Errorf("no sources provided to sort")
	}
	mediaPaths := map[string]struct{}{}
	for _, src := range source { // Make paths absolute
		cleansrc := cxt.AbsPath(src)
		if strings.HasPrefix(cleansrc, cxt.Root) {
			return fmt.Errorf("sorting source directory cannot be from within project")
		}
		info, err := os.Stat(cleansrc)
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() { // Add single files
			mediaPaths[cleansrc] = struct{}{}
		} else if info.IsDir() {
			mediaItems, err := format.GetMediaFromDirectory(cleansrc)
			if err != nil {
				return err
			}
			for _, media := range mediaItems { // Add files from directory
				mediaPaths[media.Path] = struct{}{}
			}
		}
	}

	if len(mediaPaths) == 0 {
		return nil // Nothing to do...
	}

	// Ensure sorted dir exists
	sortedDir := filepath.Join(cxt.Root, filepath.FromSlash(cxt.Config.Sorted))
	if err := os.Mkdir(sortedDir, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	// Move files into their folders
	for sourcePath := range mediaPaths {
		date, err := GetMediaDate(sourcePath)
		if err != nil {
			return err
		}
		folderPath := filepath.Join(sortedDir, FormatDate(date))
		if err = os.Mkdir(folderPath, 0755); err != nil && !os.IsExist(err) {
			return err
		}
		destPath := UniqueName(filepath.Join(folderPath, filepath.Base(sourcePath)))
		if copyFiles {
			log.Println("Copying:", sourcePath, "--->", destPath)
			if err = <-copy.File(sourcePath, destPath); err != nil {
				return err
			}
		} else {
			log.Println("Moving:", sourcePath, "--->", destPath)
			if err = os.Rename(sourcePath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}
