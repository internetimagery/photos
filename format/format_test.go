package format

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewMedia(t *testing.T) {
	event := "18-12-08 event"

	// Test filename with tags
	test1 := event + "_002[one-two three].jpg"
	media1 := NewMedia(test1)
	if media1.Event != event || media1.Index != 2 || media1.Path != test1 || media1.Ext != "jpg" || len(media1.Tags) != 2 {
		fmt.Println("Failed on", test1)
		fmt.Println(media1)
		t.Fail()
	}

	// Test filename without tags
	test2 := event + "_202.png"
	media2 := NewMedia(test2)
	if media2.Event != event || media2.Index != 202 || media2.Path != test2 || media2.Ext != "png" || len(media2.Tags) != 0 {
		fmt.Println("Failed on", test2)
		fmt.Println(media2)
		t.Fail()
	}

	// Test filename unformatted
	test3 := "my_fav_picture.jpeg"
	media3 := NewMedia(test3)
	if media3.Index != 0 {
		fmt.Println("Failed on", test3)
		fmt.Println(media3)
		t.Fail()
	}

}

type testCase struct {
	Value string
	Media Media
}

func TestFormatName(t *testing.T) {
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
			fmt.Println(err)
			t.Fail()
		} else if result != expect.Value {
			fmt.Printf("Expected '%s', got '%s'\n", expect.Value, result)
			t.Fail()
		}
	}
}

func TestGetMediaFromDirectory(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestGetMediaFromDirectory")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	rootName := "18-05-12 event"
	rootPath := filepath.Join(tmpDir, rootName)
	testFiles := map[string]*Media{
		filepath.Join(rootPath, "18-05-12 event_034.img"):                &Media{Event: "18-05-12 event", Index: 34, Ext: "img"},
		filepath.Join(rootPath, "18-05-12 event_034[one two-three].img"): &Media{Event: "18-05-12 event", Index: 34, Tags: []string{"one", "two-three"}, Ext: "img"},
		filepath.Join(rootPath, "12-10-12 event_034.png"):                &Media{Event: "18-05-12 event", Ext: "png"},
		filepath.Join(rootPath, "document_scanned.jpg"):                  &Media{Event: "18-05-12 event", Ext: "jpg"},
		filepath.Join(rootPath, TEMPPREFIX+"document_scanned.jpg"):       nil,
	}
	err = os.Mkdir(rootPath, 0755)
	if err != nil {
		t.Fatal(err)
	}
	for path := range testFiles {
		err = ioutil.WriteFile(path, []byte("testing123"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
	result, err := GetMediaFromDirectory(rootPath)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(result) != 4 {
		fmt.Printf("Expected 4 media items. Got %d\n", len(result))
		t.Fail()
	}
	for _, test := range result {
		if test == nil {
			fmt.Println("Should not have picked up", test.Path)
		} else {
			expect := testFiles[test.Path]
			if test.Event != expect.Event || test.Ext != expect.Ext || len(test.Tags) != len(expect.Tags) {
				fmt.Println("Test failed at", test.Path)
				fmt.Println(test)
				t.Fail()
			}
		}
	}

}
