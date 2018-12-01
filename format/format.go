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

// Format parts
var dateFmt = `\d{2}\-\d{2}\-\d{2}`
var dateLayout = `06-01-02`
var eventNameFmt = `[\w\- &@!%#]+` // Valid event
var indexFmt = `\d+`
var versionFmt = `\d+`
var TagsFmt = `[\w\-_ ]+` // TagsFmt : format required to be a valid tag
var extFmt = `[a-zA-Z0-9]+`

// Format example:
// Date: 18-10-10, Event: myevent, Index: 20, Version 10, Tags: one two, Ext: txt
// Format: 18-10-10 myevent_20.10[one two].txt
var mediaReg = regexp.MustCompile(fmt.Sprintf(`^(%s) (%s)_(%s)(?:\.(%s))?(?:\[(%s)\])?\.(%s)$`,
	dateFmt,
	eventNameFmt,
	indexFmt,
	versionFmt,
	TagsFmt,
	extFmt))
var eventReg = regexp.MustCompile(fmt.Sprintf(`^(?:(%s) )?(%s)$`,
	dateFmt,
	eventNameFmt))

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
		if date, err := time.ParseInLocation(dateLayout, parts[1], time.Local); err == nil {
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
	Path    string              // File name
	Date    *time.Time          // Date
	Event   string              // Event name (parent folder)
	Index   int                 // ID of media
	Version int                 // Current version of media
	Tags    map[string]struct{} // Any Tags
	Ext     string              // Extension / file type
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
		if version, err := strconv.Atoi(parts[4]); err == nil {
			media.Version = version
		}

		if date, err := time.ParseInLocation(dateLayout, parts[1], time.Local); err == nil {
			media.Date = &date
		} else {
			media.Index = 0 // Mark invalid
		}

		if len(parts[5]) > 0 {
			for _, tagname := range strings.Split(parts[5], " ") {
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
	if !regexp.MustCompile("^"+eventNameFmt+"$").MatchString(media.Event) || strings.TrimSpace(media.Event) == "" {
		return "", fmt.Errorf("Bad Event: '%s'", media.Event)
	}
	if media.Index <= 0 {
		return "", fmt.Errorf("Index value too low: '%d'", media.Index)
	}
	if media.Version < 0 {
		return "", fmt.Errorf("Cannot have negative version '%d'", media.Version)
	}
	if !regexp.MustCompile("^"+extFmt+"$").MatchString(media.Ext) || strings.TrimSpace(media.Ext) == "" {
		return "", fmt.Errorf("Bad extension: '%s'", media.Ext)
	}
	tagTest := regexp.MustCompile("^" + TagsFmt + "$")
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
	version := ""
	if media.Version > 0 {
		version = "." + strconv.Itoa(media.Version)
	}
	return fmt.Sprintf("%s %s_%03d%s%s.%s", media.Date.Format(dateLayout), media.Event, media.Index, version, tags, ext), nil
}
