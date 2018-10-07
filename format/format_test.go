package format

import (
	"fmt"
	"testing"
)

func TestRegexFromEvent(t *testing.T) {
	event := "18-12-08 event"
	tests := map[string]bool{
		"18-12-08 event_0032[someone something].jpg": true,
		"18-12-08 event_0032.png":                    true,
		"18-12-08 event023.png":                      false,
		"_023.png":                                   false,
		"18-12-08 event_0032[a thing.png":            false,
	}

	reg := RegFromEvent(event)
	for test, match := range tests {
		if match != reg.MatchString(test) {
			fmt.Println("Failed on", test)
			t.Fail()
		}
	}
}
