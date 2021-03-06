package format

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// TEMPPREFIX : Prefix for temporary working files. Ignore these files.
var TEMPPREFIX = `tmp-` // Prefix for temporary working files
// EventReg : Event name. Restrictive characters
var EventReg = `[\w\-_ ]+` // Valid event
// IndexReg : Valid index
var IndexReg = `\d+` // Valid Index
// TagReg : Valid tag characters
var TagReg = `[\w\-_ ]+` // Valid Tags
var extReg = `\w+`
var formatReg = regexp.MustCompile(fmt.Sprintf(`(%s)_(%s)(?:\[(%s)\])?\.(%s)$`, EventReg, IndexReg, TagReg, extReg))

// MakeTempPath : Apply temporary prefix to filepath
func MakeTempPath(path string) string {
	dirname := filepath.Dir(path)
	basename := filepath.Base(path)
	return filepath.Join(dirname, TEMPPREFIX+basename)
}

// IsTempPath : The counterpart to MakeTempPath. Detect if given path is indeed a temp path
func IsTempPath(path string) bool {
	basename := filepath.Base(path)
	return strings.HasPrefix(basename, TEMPPREFIX)
}

// IsUsable : Helper function that determines if a path should be considered usable as media
func IsUsable(path string) bool {
	return !IsTempPath(path) && // Do not want temp paths
		filepath.Base(path)[0] != '.' && // Cannot be a file starting with .
		!strings.HasSuffix(path, ".yaml") // Cannot be a config file
}

// Media : Container for information about media item
type Media struct {
	Path  string              // File name
	Event string              // Event name (parent folder)
	Index int                 // ID of media
	Tags  map[string]struct{} // Any Tags
	Ext   string              // Extension / file type
}

// NewMedia : Create new media representation
func NewMedia(filename string) *Media {
	media := new(Media)
	media.Path = filename
	ext := filepath.Ext(filename)
	if ext != "" {
		media.Ext = filepath.Ext(filename)[1:]
	}
	media.Tags = make(map[string]struct{})
	parts := formatReg.FindStringSubmatch(filename)
	if len(parts) > 0 {
		media.Event = parts[1]
		index, _ := strconv.Atoi(parts[2])
		media.Index = index
		if len(parts[3]) > 0 {
			for _, tagname := range strings.Split(parts[3], " ") {
				media.Tags[tagname] = struct{}{}
			}
		}
	}
	return media
}

// FormatName : Given the current settings (which may have been modified), validate and format a corresponding name.
func (media *Media) FormatName() (string, error) {
	// Validate our inputs
	if !regexp.MustCompile("^"+EventReg+"$").MatchString(media.Event) || strings.TrimSpace(media.Event) == "" {
		return "", fmt.Errorf("Bad Event: '%s'", media.Event)
	}
	if media.Index <= 0 {
		return "", fmt.Errorf("Index value too low: '%d'", media.Index)
	}
	if !regexp.MustCompile("^"+extReg+"$").MatchString(media.Ext) || strings.TrimSpace(media.Ext) == "" {
		return "", fmt.Errorf("Bad extension: '%s'", media.Ext)
	}
	tagTest := regexp.MustCompile("^" + TagReg + "$")
	for tag := range media.Tags {
		if !tagTest.MatchString(tag) || strings.TrimSpace(tag) == "" {
			return "", fmt.Errorf("Bad tag: '%s'", tag)
		}
	}

	tags := ""
	if len(media.Tags) > 0 {
		tagnames := []string{}
		for tagname := range media.Tags {
			tagnames = append(tagnames, tagname)
		}
		sort.Strings(tagnames)
		tags = fmt.Sprintf("[%s]", strings.Join(tagnames, " "))
	}
	ext := strings.ToLower(media.Ext)
	return fmt.Sprintf("%s_%03d%s.%s", media.Event, media.Index, tags, ext), nil
}

// GetMediaFromDirectory : Walk through directory, and return a list of media items represented there
func GetMediaFromDirectory(dirPath string) ([]*Media, error) {
	mediaList := []*Media{}
	files, err := ioutil.ReadDir(dirPath)
	event := filepath.Base(dirPath)
	if err != nil {
		return mediaList, err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	for _, file := range files {
		if file.Mode().IsRegular() && IsUsable(filepath.Join(dirPath, file.Name())) { // Ignore unusable files
			fullPath := filepath.Join(dirPath, file.Name())
			media := NewMedia(fullPath)
			if media.Event != event {
				media.Index = 0 // Index 0 means unformatted
				media.Event = event
			}
			mediaList = append(mediaList, media)
		}
	}
	return mediaList, nil
}
