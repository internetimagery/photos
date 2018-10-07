package format

import (
	"fmt"
	"regexp"
)

var eventReg = `(%s)`
var indexReg = `_(\d+)`
var tagReg = `[\w\-_ ]+`
var tagsReg = fmt.Sprintf(`(?:\[(%s)\])?`, tagReg)
var extReg = `\.(\w+)`

// Media : Container for information about media item
type Media struct {
	Name  string   // File name
	Event string   // Event name (parent folder)
	Index int      // ID of media
	Tags  []string // Any Tags
	Ext   string   // Extension / file type
}

// RegFromEvent : Form regex from event name
func RegFromEvent(event string) *regexp.Regexp {
	return regexp.MustCompile(
		fmt.Sprintf(eventReg, regexp.QuoteMeta(event)) +
			indexReg + tagsReg + extReg)
}
