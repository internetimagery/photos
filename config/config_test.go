// Testing configuration creation, loading and searching
package config

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/testutil"
)

func TestNewConfig(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	handle := new(bytes.Buffer)
	conf := NewConfig("test") // Create new config data
	err := conf.Save(handle)
	if err != nil {
		tu.Fail(err)
	}
	testData := handle.Bytes()
	verifyStruct := make(map[string]interface{}) // Load config for basic test
	err = json.Unmarshal(testData, &verifyStruct)
	if err != nil {
		tu.Fail(err)
	}

	// Verify basic groups are present
	if _, ok := verifyStruct["compress"]; !ok {
		tu.Fail("Config missing compress group")
	}

	if _, ok := verifyStruct["backup"]; !ok {
		tu.Fail("Config missing backup group")
	}

	if _, ok := verifyStruct["id"]; !ok {
		tu.Fail("Config missing id group")
	}
}

func TestCompressCommand(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	testData := `
	{
	 "compress":[
	    ["*.jpg *.png", "image"],
			["*.mp4 video/*", "video"],
			["path/to/*/*", "path"],
	    ["*", "all"]
	 ]
 }`
	handle := bytes.NewReader([]byte(testData))
	conf, err := LoadConfig(handle)
	if err != nil {
		tu.Fail(err)
	}
	// Do some testing!
	tests := map[string]string{
		"anything":                         "all",
		"somepic.JPG":                      "image",
		"path/to/other/pic.png":            "image", // Ignoring paths
		"video.mp4":                        "video",
		filepath.Join("video", "file.vid"): "all",
	}
	for test, expect := range tests {
		command := conf.Compress.GetCommand(test)
		if command != expect {
			tu.FailE(expect, command)
		}
	}
}

func TestBackupCommand(t *testing.T) {
	tu := testutil.NewTestUtil(t)

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
		tu.Fail(err)
	}
	// Do some testing!
	tests := map[string]map[string]bool{
		"local":   map[string]bool{"local": true},
		"remote*": map[string]bool{"dropbox": true, "amazon": true},
	}
	for test, expect := range tests {
		commands := conf.Backup.GetCommands(test)
		if len(commands) == 0 {
			tu.Fail("No commands returned for", test)
		}
		for _, command := range commands {
			if !expect[command] {
				tu.FailE(command, test)
			}
		}
	}
}
