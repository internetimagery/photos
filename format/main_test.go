// Testing formatting
package format

import "testing"

func TestMatch(t *testing.T) {
	good1 := Match("dir_name", "dir_name_001.mov")
	good2 := Match("dir_name", "dir_name_001[stuff more stuff].jpg")
	bad1 := Match("dir_name", "whatever.png")
	bad2 := Match("dir_name", "whatever_002.png")
	bad3 := Match("dir_name", "dir_name[tags].png")
	bad4 := Match("dir_name", "dir_name_004{why not}.png")

	if !good1.Formatted {
		t.Fail()
	}
	if !good2.Formatted {
		t.Fail()
	}
	if bad1.Formatted {
		t.Fail()
	}
	if bad2.Formatted {
		t.Fail()
	}
	if bad3.Formatted {
		t.Fail()
	}
	if bad4.Formatted {
		t.Fail()
	}
}
