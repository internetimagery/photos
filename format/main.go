// Format functionality
package format

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Media struct {
	Name      string
	Formatted bool
	Index     int
	Tags      []string
	Ext       string
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
				media.Tags = append(media.Tags, tag)
			}
		}
	}
	return media, nil
}

func Match(dir string, names []string) ([]*Media, error) {
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
