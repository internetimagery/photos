package sort

import (
	"os"
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
