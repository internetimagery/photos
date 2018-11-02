package sort

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/internetimagery/photos/context"
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

	// Try processing exif data
	if exifData, err := exif.Decode(handle); err == nil {
		taken, err := exifData.DateTime()
		if err == nil {
			return taken, nil
		}
		return time.Time{}, err
	}
	info, err := handle.Stat()
	if err == nil {
		return info.ModTime(), nil
	}
	return time.Time{}, err
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
func SortMedia(cxt *context.Context) error {
	infos, err := ioutil.ReadDir(cxt.WorkingDir)
	if err != nil {
		return err
	}

	// Move files into their folders
	for _, info := range infos {
		if !info.IsDir() {
			sourcePath := filepath.Join(cxt.WorkingDir, info.Name())
			date, err := GetMediaDate(sourcePath)
			if err != nil {
				return err
			}
			folderPath := filepath.Join(cxt.WorkingDir, FormatDate(date))
			if err = os.Mkdir(folderPath, 0755); err != nil && !os.IsExist(err) {
				return err
			}
			fmt.Println("destPath", folderPath, info.Name())
			destPath := UniqueName(filepath.Join(folderPath, info.Name()))
			log.Println("Moving:", sourcePath, "--->", destPath)
			if err = os.Rename(sourcePath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}
