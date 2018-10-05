// Testing configuration creation, loading and searching
package config

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestLoadConfig(t *testing.T) {
	testData := `
	{
	 "compress":[
	    {"name": "filter1 filter2", "command": "command1"},
	    {"name": "filter3", "command": "command2"}
	 ],
	 "backup":[
	    {"name": "optionA", "command": "command3"},
	    {"name": "optionB", "command": "command4"}
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
