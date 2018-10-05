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

func TestBackupCommand(t *testing.T) {
	testData := `
	{
	 "backup":[
	    ["remote-dropbox", "dropbox"],
			["remote-amazon", "amazon"],
	    ["local", "local"]
	 ]
 }`
	handle := bytes.NewReader([]byte(testData))
	conf, err := LoadConfig(handle)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	// Do some testing!
	tests := map[string]map[string]bool{
		"local":   map[string]bool{"local": true},
		"remote*": map[string]bool{"dropbox": true, "amazon": true},
	}
	for test, expect := range tests {
		commands := conf.Backup.GetCommands(test)
		if len(commands) == 0 {
			fmt.Println("No commands returned for", test)
			t.Fail()
		}
		for _, command := range commands {
			if !expect[command] {
				fmt.Printf("Got '%s' while testing '%s'\n", command, test)
				t.Fail()
			}
		}
	}
}
