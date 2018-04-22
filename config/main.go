// Repo configuration
package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/rs/xid"
)

// Config data
type Config struct {
	ID       string            `json:"id"`       // Unique ID
	Name     string            `json:"name"`     // Optional name
	Root     string            `json:"-"`        // Location of the config file (not stored in config)
	Commands map[string]string `json:"commands"` // Command names and paths
}

// Create a new config file
func NewConfig() *Config {
	// Perform initial setup here
	conf := new(Config)
	conf.ID = xid.New().String() // Generate random ID
	return conf
}

// Load config from file
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := new(Config)
	err = json.Unmarshal(data, conf)
	if conf.ID == "" { // ID is a required field. Error if not found.
		return nil, errors.New("Missing ID")
	}
	conf.Root = path
	return conf, err
}

// Store config data to file
func (self Config) Save(path string) error {
	// Check id first
	conf, err := LoadConfig(path)
	if err == nil && conf.ID != self.ID {
		return errors.New("ID's do not match. Can't save config.")
	}

	data, err := json.MarshalIndent(self, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 664)
	return err
}
