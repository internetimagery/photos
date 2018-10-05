// Generate or go hunting for a configuration file.
package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Config : Base class to access root configuration
type Config struct {
	Compress map[string]string `json:"compress"` // Compression commands
	Backup   map[string]string `json:"backup"`   // Backup commands
}

// NewConfig build barebones data to get started on a new config file
func NewConfig(writer io.Writer) error {
	newConfig := new(Config) // Create empty config, and add some default info to assist in fleshing out properly
	newConfig.Compress = make(map[string]string)
	newConfig.Compress["*.example2 *.example2"] = "// command to run on files ending with '.example1' or '.example2'"
	newConfig.Backup = make(map[string]string)
	newConfig.Backup["placeofbackup"] = "// command to run when selecting this backup option 'placeofbackup'"
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
