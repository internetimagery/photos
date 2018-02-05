// Testing formatting
package format

import (
	"testing"
)

func TestMatch(t *testing.T) {
	good1, _ := Match("dir_name", []string{"dir_name_001.mov"})
	good2, _ := Match("dir_name", []string{"dir_name_001[stuff more stuff].jpg"})
	bad1, _ := Match("dir_name", []string{"whatever.png"})
	bad2, _ := Match("dir_name", []string{"whatever_002.png"})
	bad3, _ := Match("dir_name", []string{"dir_name[tags].png"})
	bad4, _ := Match("dir_name", []string{"dir_name_004{why not}.png"})
	bad5, _ := Match("dir_name", []string{"dir_name_005[why yes.png"})

	if !good1[0].Formatted || good1[0].Index != 1 {
		t.Fail()
	}
	if !good2[0].Formatted || good2[0].Index != 1 || good2[0].Tags[0] != "stuff" {
		t.Fail()
	}
	if bad1[0].Formatted || bad2[0].Formatted || bad3[0].Formatted || bad4[0].Formatted || bad5[0].Formatted {
		t.Fail()
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
