package format

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/testutil"
)

func TestNewMedia(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	event := "18-12-08 event"

	// Test filename with tags
	test1 := event + "_002[one-two three].jpg"
	media1 := NewMedia(test1)
	if media1.Event != event || media1.Index != 2 || media1.Path != test1 || media1.Ext != "jpg" || len(media1.Tags) != 2 {
		tu.Fail("Failed on", test1, media1)
	}

	// Test filename without tags
	test2 := event + "_202.png"
	media2 := NewMedia(test2)
	if media2.Event != event || media2.Index != 202 || media2.Path != test2 || media2.Ext != "png" || len(media2.Tags) != 0 {
		tu.Fail("Failed on", test2, media2)
	}

	// Test filename unformatted
	test3 := "my_fav_picture.jpeg"
	media3 := NewMedia(test3)
	if media3.Index != 0 {
		tu.Fail("Failed on", test3, media3)
	}

}

type testCase struct {
	Value string
	Media Media
}

func TestFormatName(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	tests := []testCase{
		testCase{"event01_020.png", Media{Event: "event01", Index: 20, Ext: "png"}},
		testCase{"18-12-07 event_1234[one two].jpeg", Media{Event: "18-12-07 event", Index: 1234, Tags: []string{"one", "two"}, Ext: "jpeg"}},
		testCase{"", Media{Event: "some event/event", Index: 2, Ext: "jpg"}},
		testCase{"", Media{Event: "evento", Index: -1, Ext: "png"}},
		testCase{"", Media{Event: "eventing", Index: 23, Ext: "$$$"}},
		testCase{"", Media{Event: "  ", Index: 23, Ext: "thing"}},
		testCase{"", Media{Event: "eventer", Index: 12, Ext: ""}},
	}

	for _, expect := range tests {
		result, err := expect.Media.FormatName()
		if err != nil && expect.Value != "" {
			tu.Fail(err)
		} else if result != expect.Value {
			tu.FailE(expect.Value, result)
		}
	}
}

func TestGetMediaFromDirectory(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.TempDir("TestGetMediaFromDirectory")()

	rootName := "18-05-12 event"
	rootPath := filepath.Join(tu.Dir, rootName)
	testFiles := map[string]*Media{
		filepath.Join(rootPath, "18-05-12 event_034.img"):                &Media{Event: "18-05-12 event", Index: 34, Ext: "img"},
		filepath.Join(rootPath, "18-05-12 event_034[one two-three].img"): &Media{Event: "18-05-12 event", Index: 34, Tags: []string{"one", "two-three"}, Ext: "img"},
		filepath.Join(rootPath, "12-10-12 event_034.png"):                &Media{Event: "18-05-12 event", Ext: "png"},
		filepath.Join(rootPath, "document_scanned.jpg"):                  &Media{Event: "18-05-12 event", Ext: "jpg"},
		filepath.Join(rootPath, TEMPPREFIX+"document_scanned.jpg"):       nil,
	}
	err := os.Mkdir(rootPath, 0755)
	if err != nil {
		tu.Fatal(err)
	}
	for path := range testFiles {
		tu.NewFile(path, "")
	}
	result, err := GetMediaFromDirectory(rootPath)
	if err != nil {
		tu.Fail(err)
	}
	if len(result) != 4 {
		tu.Fail("Expected 4 media items. Got", len(result))
	}
	for _, test := range result {
		if test == nil {
			t.Log("Should not have picked up", test.Path)
		} else {
			expect := testFiles[test.Path]
			if test.Event != expect.Event || test.Ext != expect.Ext || len(test.Tags) != len(expect.Tags) {
				tu.Fail("Test failed at", test.Path)
			}
		}
	}

}
