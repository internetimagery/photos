// Format functionality
package format

import (
	"regexp"
	"strconv"
	"strings"
)

type Media struct {
	Name      string
	Formatted bool
	Index     int
	Tags      []string
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
	if len(parts) > 0 {
		media.Formatted = true
		index, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		media.Index = index
		media.Tags = strings.Split(parts[2], " ")
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
