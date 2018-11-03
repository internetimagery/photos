package format

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/internetimagery/photos/testutil"
)

func TestTempPath(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	testpath := "/one/two/three/four.five"
	temppath := MakeTempPath(testpath)
	if !IsTempPath(temppath) {
		tu.Fail("Failed to match temp path", temppath)
	}
	if IsTempPath(testpath) {
		tu.Fail("False positive on temp path", testpath)
	}
}

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

func TestFormatName(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	type testCase struct {
		Value string
		Media Media
	}

	tests := []testCase{
		testCase{"event01_020.png", Media{Event: "event01", Index: 20, Ext: "png"}},
		testCase{"18-12-07 event_1234[one two].jpeg", Media{Event: "18-12-07 event", Index: 1234, Tags: map[string]struct{}{"one": struct{}{}, "two": struct{}{}}, Ext: "jpeg"}},
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
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "18-05-12 event")
	result, err := GetMediaFromDirectory(event)
	if err != nil {
		tu.Fail(err)
	}

	if len(result) != 4 {
		tu.Fail("Expected 4 media items. Got", len(result))
	}

	testFiles := map[string]Media{
		filepath.Join(event, "18-05-12 event_034.img"):                Media{Event: "18-05-12 event", Index: 34, Ext: "img"},
		filepath.Join(event, "18-05-12 event_034[one two-three].img"): Media{Event: "18-05-12 event", Index: 34, Tags: map[string]struct{}{"one": struct{}{}, "two-three": struct{}{}}, Ext: "img"},
		filepath.Join(event, "12-10-12 event_034.png"):                Media{Event: "18-05-12 event", Ext: "png"},
		filepath.Join(event, "document_scanned.jpg"):                  Media{Event: "18-05-12 event", Ext: "jpg"},
	}

	for _, testFile := range result {
		if strings.HasPrefix(testFile.Path, ".") {
			tu.Fail("Caught file beginning with '.':", testFile.Path)
		}
		tu.AssertExists(testFile.Path)
		expect := testFiles[testFile.Path]
		if testFile.Event != expect.Event || testFile.Ext != expect.Ext || len(testFile.Tags) != len(expect.Tags) {
			tu.Fail("Test failed at", testFile.Path)
		}
	}

}
