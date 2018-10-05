// Generate or go hunting for a configuration file.
package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

///////////////////////////////////
// TODO: Adding filter to all commands means we can use combined filter functionality
// TODO: ie: We can apply filenames to the filter on compress commands to get the first command that matches.
// TODO: ie: We can apply the same to match backup commands to a vague command request. To run more than one command at once.
// TODO: Should config be able to run commands when requested? Utility function?
// TODO: https://golang.org/pkg/path/filepath/#Match split filename on space? use go-shlex?
// TODO: Should config be able to detect which command to run with filtering Compress?
// TODO: Do we even need a "loadConfig" test? type checking and struct filling sort that out...
///////////////////////////////////

// Command : Structure for a command
type Command struct {
	Name    string `json:"name"`    // Name for command. Can be used as a filter for compress matches etc
	Command string `json:"command"` // Command to run
}

// Category : Groups commands together. Facilitates finding commands by name
type Category []Command

// Config : Base class to access root configuration
type Config struct {
	Compress Category `json:"compress"` // Compression commands
	Backup   Category `json:"backup"`   // Backup commands
}

// NewConfig build barebones data to get started on a new config file
func NewConfig(writer io.Writer) error {
	newConfig := new(Config) // Create empty config, and add some default info to assist in fleshing out properly
	newConfig.Compress = Category{
		Command{"*.example2 *.example2", "// command to run on files ending with '.example1' or '.example2'"}}
	newConfig.Backup = Category{
		Command{"placeofbackup", "// command to run when selecting this backup option 'placeofbackup'"}}
	newData, err := json.Marshal(newConfig)
	if err != nil {
		return err
	}
	_, err = writer.Write(newData)
	return err
}

// LoadConfig : Load and populate a new Config from existing config data
func LoadConfig(reader io.Reader) (*Config, error) {
	loadedData, err := ioutil.ReadAll(reader) // Load the data to process
	if err != nil {
		return nil, err
	}
	loadedConfig := new(Config)
	err = json.Unmarshal(loadedData, loadedConfig)
	return loadedConfig, err
}
