package format

import (
	"fmt"
	"testing"
)

func TestNewMedia(t *testing.T) {
	event := "18-12-08 event"

	// Test filename with tags
	test1 := event + "_002[one-two three].jpg"
	media1 := NewMedia(test1)
	if media1.Event != event || media1.Index != 2 || media1.Name != test1 || media1.Ext != "jpg" || len(media1.Tags) != 2 {
		fmt.Println("Failed on", test1)
		fmt.Println(media1)
		t.Fail()
	}

	// Test filename without tags
	test2 := event + "_202.png"
	media2 := NewMedia(test2)
	if media2.Event != event || media2.Index != 202 || media2.Name != test2 || media2.Ext != "png" || len(media2.Tags) != 0 {
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

func TestFormatName(t *testing.T) {
	tests := make(map[string]Media)

	tests["event01_020.png"] = Media{Event: "event01", Index: 20, Ext: "png"}
	tests["18-12-07 event_1234[one two].jpeg"] = Media{Event: "18-12-07 event", Index: 1234, Tags: []string{"one", "two"}, Ext: "jpeg"}
	tests[""] = Media{Event: "some event/event", Index: 0, Ext: "$$$"}

	for expect, test := range tests {
		result, err := test.FormatName()
		if err != nil && expect != "" {
			fmt.Println(err)
			t.Fail()
		} else if result != expect {
			fmt.Printf("Expected '%s', got '%s'\n", expect, result)
			t.Fail()
		}
	}
}
