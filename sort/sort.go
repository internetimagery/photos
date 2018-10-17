package sort

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GetMediaDate : Get modification date, or date taken (EXIF data) from file
func GetMediaDate(filePath string) (time.Time, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
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
	name := filename[:len(ext)]
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
