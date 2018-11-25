package format

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

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

	// Test filename with tags
	test := "18-12-08 event_002[one-two three].jpg"
	media := NewMedia(test)
	if media.Date.Equal(time.Date(2018, 12, 8, 0, 0, 0, 0, time.Local)) || media.Event != "event" || media.Index != 2 || media.Path != test || media.Ext != "jpg" || len(media.Tags) != 2 {
		tu.Fail("Failed on", test, media)
	}

	// Test filename without tags
	test = "18-12-08 event_202.png"
	media = NewMedia(test)
	if media.Date.Equal(time.Date(2018, 12, 8, 0, 0, 0, 0, time.Local)) || media.Event != "event" || media.Index != 202 || media.Path != test || media.Ext != "png" || len(media.Tags) != 0 {
		tu.Fail("Failed on", test, media)
	}

	// Test filename unformatted
	test = "my_fav_picture.jpeg"
	media = NewMedia(test)
	if media.Index != 0 {
		tu.Fail("Failed on", test, media)
	}

	// Test filename no extension
	test = "18-12-08 event_101"
	media = NewMedia(test)
	if media.Index != 0 || media.Ext != "" {
		tu.Fail("Failed on", test, media)
	}

	// Test filename bad date
	test = "20-30-40 event_101.log"
	media = NewMedia(test)
	if media.Index != 0 || media.Ext != "" {
		tu.Fail("Failed on", test, media)
	}

}

func TestFormatName(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	type testCase struct {
		Value string
		Media Media
	}

	testTime := time.Date(2018, 12, 7, 0, 0, 0, 0, time.Local)
	tests := []testCase{
		testCase{"18-12-07 event01_020.png", Media{Event: "event01", Index: 20, Ext: "png", Date: &testTime}},
		testCase{"18-12-07 event_1234[one two].jpeg", Media{Event: "event", Index: 1234, Tags: map[string]struct{}{"one": struct{}{}, "two": struct{}{}}, Ext: "jpeg", Date: &testTime}},
		testCase{"", Media{Event: "some event/event", Index: 2, Ext: "jpg", Date: &testTime}},
		testCase{"", Media{Event: "evento", Index: -1, Ext: "png", Date: &testTime}},
		testCase{"", Media{Event: "eventing", Index: 23, Ext: "$$$", Date: &testTime}},
		testCase{"", Media{Event: "  ", Index: 23, Ext: "thing", Date: &testTime}},
		testCase{"", Media{Event: "eventer", Index: 12, Ext: "", Date: &testTime}},
		testCase{"", Media{Event: "event", Index: 2, Ext: "jpg"}},
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
	result := tu.Must(GetMediaFromDirectory(event)).([]*Media)

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

func TestIsUsable(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	testpath1 := "/one/two/three.four"     // normal file
	testpath2 := "/one/two/three.yaml"     // config file
	testpath3 := "/one/two/.three"         // dotted file
	testpath4 := "/one/two/tmp-three.four" // temp file

	if !IsUsable(testpath1) {
		tu.Fail("Failed on normal path")
	}
	if IsUsable(testpath2) || IsUsable(testpath3) || IsUsable(testpath4) {
		tu.Fail("Failed on unusable path")
	}
}
