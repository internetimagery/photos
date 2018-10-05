// Testing configuration creation, loading and searching
package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
)

func TestNewConfig(t *testing.T) {
	handle := new(bytes.Buffer)
	err := NewConfig(handle) // Create new config data
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	testData := handle.Bytes()
	verifyStruct := make(map[string]interface{}) // Load config for basic test
	err = json.Unmarshal(testData, &verifyStruct)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	// Verify basic groups are present
	if _, ok := verifyStruct["compress"]; !ok {
		fmt.Println("Config missing compress group")
		t.Fail()
	}

	if _, ok := verifyStruct["backup"]; !ok {
		fmt.Println("Config missing backup group")
		t.Fail()
	}
}

func TestCompressCommand(t *testing.T) {
	testData := `
	{
	 "compress":[
	    ["*.jpg *.png", "image"],
			["*.mp4 video/*", "video"],
	    ["*", "all"]
	 ]
 }`
	handle := bytes.NewReader([]byte(testData))
	conf, err := LoadConfig(handle)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	// Do some testing!
	tests := map[string]string{
		"anything":                         "all",
		"somepic.JPG":                      "image",
		"other/pic.png":                    "image",
		"video.mp4":                        "video",
		filepath.Join("video", "file.vid"): "video",
	}
	for test, expect := range tests {
		command := conf.Compress.GetCommand(test)
		if command != expect {
			fmt.Printf("Expected '%s' but got '%s' while testing '%s'\n", expect, command, test)
			t.Fail()
		}
	}
}

func TestLoadConfig(t *testing.T) {
	testData := `
	{
	 "compress":[
	    ["filter1 filter2", "command1"],
	    ["filter3", "command2"]
	 ],
	 "backup":[
	    ["optionA", "command3"],
	    ["optionB", "command4"]
	 ]
	}`

	// Load our mock data
	handle := bytes.NewReader([]byte(testData))
	conf, err := LoadConfig(handle)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if conf == nil {
		fmt.Println("Huh? Missing config?")
		t.Fail()
	}
}
