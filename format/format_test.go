// Testing formatting
package format

import (
	"fmt"
	"testing"
)

func TestSanitize(t *testing.T) {
	tests := make(map[string]string)
	tests["one.two"] = "one.two"
	tests["three/four"] = "three_four"
	tests["Hello, 世界"] = "Hello_ 世界"

	for test, expect := range tests {
		result := Sanitize(test)
		if result != expect {
			fmt.Println("Testing:", test, "Expected:", expect, "Got:", result)
			t.Fail()
		}
	}
}

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
		"",
	})
	_, err := Match(" ", []string{"fail.me"})

	// Test for empty directory
	if err == nil {
		t.Fail()
	}
	// Test for well formatted without tags
	if !good[0].Formatted || good[0].Index != 1 || good[0].Ext != ".mov" {
		t.Fail()
	}
	// Test for well formatted with tags
	if !good[1].Formatted || good[1].Index != 1 || good[1].Tags[0] != "stuff" {
		t.Fail()
	}
	// Test for no repeated tags
	if len(good[1].Tags) != 2 {
		t.Fail()
	}
	// Test for different types of bad formats
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
