package format

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TEMPPREFIX : Prefix for temporary working files. Ignore these files.
var TEMPPREFIX = `tmp-` // Prefix for temporary working files

// DateReg : Date
var DateReg = `\d{2}\-\d{2}\-\d{2}`

// DateLayout : Format to display date
var DateLayout = `06-01-02`

// EventNameReg : Event name. Restrictive characters
var EventNameReg = `[\w\-_ &@!%#]+` // Valid event

// IndexReg : Valid index
var IndexReg = `\d+` // Valid Index

// TagReg : Valid tag characters
var TagReg = `[\w\-_ ]+` // Valid Tags

// ExtReg : Extension
var ExtReg = `\w+`
var mediaReg = regexp.MustCompile(fmt.Sprintf(`^(%s) (%s)_(%s)(?:\[(%s)\])?\.(%s)$`, DateReg, EventNameReg, IndexReg, TagReg, ExtReg))
var eventReg = regexp.MustCompile(fmt.Sprintf(`^(?:(%s) )?(%s)$`, DateReg, EventNameReg))

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

// Event : Group of media items.
type Event struct {
	Path string     // Path to event
	Name string     // Name of event
	Date *time.Time // Time of the event (if provided)
}

// NewEvent : Create a new event object from a directory
func NewEvent(dirname string) *Event {
	event := &Event{Path: dirname}
	parts := eventReg.FindStringSubmatch(filepath.Base(dirname))
	if len(parts) > 0 {
		if date, err := time.ParseInLocation(DateLayout, parts[1], time.Local); err == nil {
			event.Date = &date
		}
		event.Name = parts[2]
	}
	return event
}

// GetMedia : Walk through directory, and return a list of media items represented there
func (event *Event) GetMedia() ([]*Media, error) {
	mediaList := []*Media{}
	if strings.TrimSpace(event.Name) == "" {
		return mediaList, fmt.Errorf("event has no name, '%s'", event.Path)
	}
	files, err := ioutil.ReadDir(event.Path)
	if err != nil {
		return mediaList, err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	for _, file := range files {
		fullPath := filepath.Join(event.Path, file.Name())
		if file.Mode().IsRegular() && IsUsable(fullPath) { // Ignore unusable files
			media := NewMedia(fullPath)
			if media.Event != event.Name {
				media.Index = 0 // Index 0 means unformatted
				media.Event = event.Name
			}
			mediaList = append(mediaList, media)
		}
	}
	return mediaList, nil
}

// Media : Container for information about media item
type Media struct {
	Path  string              // File name
	Date  *time.Time          // Date
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
		media.Ext = ext[1:]
	}
	media.Tags = make(map[string]struct{})
	parts := mediaReg.FindStringSubmatch(filepath.Base(filename))
	if len(parts) > 0 {
		media.Event = parts[2]
		index, _ := strconv.Atoi(parts[3])
		media.Index = index

		if date, err := time.ParseInLocation(DateLayout, parts[1], time.Local); err == nil {
			media.Date = &date
		} else {
			media.Index = 0 // Mark invalid
		}

		if len(parts[4]) > 0 {
			for _, tagname := range strings.Split(parts[4], " ") {
				media.Tags[tagname] = struct{}{}
			}
		}
	}
	return media
}

// FormatName : Given the current settings (which may have been modified), validate and format a corresponding name.
func (media *Media) FormatName() (string, error) {
	// Validate our inputs
	if media.Date == nil {
		return "", fmt.Errorf("missing date")
	}
	if !regexp.MustCompile("^"+EventNameReg+"$").MatchString(media.Event) || strings.TrimSpace(media.Event) == "" {
		return "", fmt.Errorf("Bad Event: '%s'", media.Event)
	}
	if media.Index <= 0 {
		return "", fmt.Errorf("Index value too low: '%d'", media.Index)
	}
	if !regexp.MustCompile("^"+ExtReg+"$").MatchString(media.Ext) || strings.TrimSpace(media.Ext) == "" {
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
	return fmt.Sprintf("%s %s_%03d%s.%s", media.Date.Format(DateLayout), media.Event, media.Index, tags, ext), nil
}
