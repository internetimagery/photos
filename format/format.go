// Format functionality
package format

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/internetimagery/photos/utility"
)

// thoughts...
// NewMediaFromFile(path)
// media object, path to associated file, name, path, id etc...
// create media object, parsing filename for metadata
// metadata can be changed to whatever. Media.Update() to rename it perhaps

type Media struct {
	Path  string   // Path to file.
	Index int      // Unique to folder. Sorting files.
	Group string   // Group media is located in. Typically the name of the containing folder.
	Type  string   // File type. Doubles as file extension.
	Hash  string   // Hash fingerprint of files content.
	Tags  []string // Extra metadata
}

// Get new media element from file.
func NewMedia(filename string) (*Media, error) {
	media := new(Media)
	media.Path = filename
	media.Type = filepath.Ext(filename)
	filehash, err := utility.GetHashFromFile(filename)
	if err != nil {
		return media, err
	}
	// TODO: regex group, index, tags
	return media, nil
}

// type Media struct {
// 	Name      string
// 	Formatted bool
// 	Index     int
// 	Tags      []string
// 	Ext       string
// }

// Remove any illegal characters, turning them into underscores.
// Allowed characters: spaces, underscores, dashes, letters, digits, fullstops, square brackets, round brackets
func Sanitize(filename string) string {
	reg := regexp.MustCompile(`[^\p{L} _\-\d\.\[\]\()]`)
	return reg.ReplaceAllString(filename, "_")
}

func (self Media) Format(dir string) string {
	name := fmt.Sprintf("%s_%03d", dir, self.Index)
	if len(self.Tags) > 0 {
		name = fmt.Sprintf("%s[%s]", name, strings.Join(self.Tags, " "))
	}
	return name + self.Ext
}

func getRegex(dir string) (*regexp.Regexp, error) {
	// Apply prefix and suffix. Return compiled regex
	suffix := "_(\\d+)(?:\\[(.+?)])?\\.(\\w+)"
	return regexp.Compile(regexp.QuoteMeta(dir) + suffix)
}

func tagIn(value string, list []string) bool {
	res := false
	for _, val := range list {
		if val == value {
			res = true
		}
	}
	return res
}

func NewMedia(regex *regexp.Regexp, name string) (*Media, error) {
	parts := regex.FindStringSubmatch(name)
	media := new(Media)
	media.Name = name
	media.Ext = filepath.Ext(name)
	if len(parts) > 0 {
		media.Formatted = true
		index, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		media.Index = index
		for _, tag := range strings.Split(parts[2], " ") {
			if tag != "" { // Skip empty tags
				if !tagIn(tag, media.Tags) {
					media.Tags = append(media.Tags, tag)
				}
			}
		}
	}
	return media, nil
}

func Match(dir string, names []string) ([]*Media, error) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return nil, errors.New("No directory specified.")
	}
	media := []*Media{}
	reg, err := getRegex(dir)
	if err != nil {
		return nil, err
	}
	for _, n := range names {
		m, err := NewMedia(reg, n)
		if err != nil {
			return nil, err
		}
		media = append(media, m)
	}
	return media, nil
}
