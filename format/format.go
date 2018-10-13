package format

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var eventReg = `[\w\-_ ]+` // Valid event
var indexReg = `\d+`       // Valid Index
var tagReg = `[\w\-_ ]+`   // Valid Tags
var extReg = `\w+`
var formatReg = regexp.MustCompile(fmt.Sprintf(`(%s)_(%s)(?:\[(%s)\])?\.(%s)$`, eventReg, indexReg, tagReg, extReg))

// Media : Container for information about media item
type Media struct {
	Path  string   // File name
	Event string   // Event name (parent folder)
	Index int      // ID of media
	Tags  []string // Any Tags
	Ext   string   // Extension / file type
}

// NewMedia : Create new media representation
func NewMedia(filename string) *Media {
	media := new(Media)
	media.Path = filename
	media.Ext = filepath.Ext(filename)[1:]
	parts := formatReg.FindStringSubmatch(filename)
	if len(parts) > 0 {
		media.Event = parts[1]
		index, _ := strconv.Atoi(parts[2])
		media.Index = index
		if len(parts[3]) > 0 {
			media.Tags = strings.Split(parts[3], " ")
		}
	}
	return media
}

// FormatName : Given the current settings (which may have been modified), validate and format a corresponding name.
func (media *Media) FormatName() (string, error) {
	// Validate our inputs
	if !regexp.MustCompile("^"+eventReg+"$").MatchString(media.Event) || strings.TrimSpace(media.Event) == "" {
		return "", fmt.Errorf("Bad Event: '%s'", media.Event)
	}
	if media.Index <= 0 {
		return "", fmt.Errorf("Index value too low: '%d'", media.Index)
	}
	if !regexp.MustCompile("^"+extReg+"$").MatchString(media.Ext) || strings.TrimSpace(media.Ext) == "" {
		return "", fmt.Errorf("Bad extension: '%s'", media.Ext)
	}
	tagTest := regexp.MustCompile("^" + tagReg + "$")
	for _, tag := range media.Tags {
		if !tagTest.MatchString(tag) || strings.TrimSpace(tag) == "" {
			return "", fmt.Errorf("Bad tag: '%s'", tag)
		}
	}

	tags := ""
	if len(media.Tags) > 0 {
		tags = fmt.Sprintf("[%s]", strings.Join(media.Tags, " "))
	}
	ext := strings.ToLower(media.Ext)
	return fmt.Sprintf("%s_%03d%s.%s", media.Event, media.Index, tags, ext), nil
}

// GetMediaFromDirectory : Walk through directory, and return a list of media items represented there
func GetMediaFromDirectory(dirPath string) ([]*Media, error) {
	mediaList := []*Media{}
	event := filepath.Base(dirPath)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return mediaList, err
	}
	for _, file := range files {
		if !file.IsDir() {
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
