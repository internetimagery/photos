// Testing formatting
package format

import (
	"testing"
)

func TestMatch(t *testing.T) {
	good, _ := Match("dir_name", []string{
		"dir_name_001.mov",
		"dir_name_001[stuff more stuff].jpg",
	})
	bad, _ := Match("dir_name", []string{
		"whatever.png",
		"whatever_002.png",
		"dir_name[tags].png",
		"dir_name_004{why not}.png",
		"dir_name_005[why yes.png",
	})

	if !good[0].Formatted || good[0].Index != 1 {
		t.Fail()
	}
	if !good[1].Formatted || good[1].Index != 1 || good[1].Tags[0] != "stuff" {
		t.Fail()
	}
	for _, b := range bad {
		if b.Formatted {
			t.Fail()
		}
	}
}

func TestFormat(t *testing.T) {
	media, _ := Match("dir_name", []string{"dir_name_003.mov"})
	name := media[0].Format("dir2_name")
	if name != "dir2_name_003.mov" {
		t.Fail()
	}
	media[0].Index = 13
	media[0].Tags = append(media[0].Tags, "some", "tag")
	media[0].Ext = ".jpg"
	name = media[0].Format("dir3_name")
	if name != "dir3_name_013[some tag].jpg" {
		t.Fail()
	}
}
