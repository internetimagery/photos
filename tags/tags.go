package tags

import "github.com/internetimagery/photos/format"

// AddTag : Apply tagnames to a file
func AddTag(file *format.Media, tagnames ...string) (*format.Media, error) {
	return file, nil
}

// RemoveTag : Remove tagnames from a file
func RemoveTag(file *format.Media, tagnames ...string) (*format.Media, error) {
	return file, nil
}
