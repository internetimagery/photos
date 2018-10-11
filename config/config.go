package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/xid"
)

// Command : Structure for a command
type Command [2]string

// CompressCategory : Groups categories together. Facilitates finding commands by filter
type CompressCategory []Command

// BackupCategory : Groups backups together. Facilitates finding commands by name
type BackupCategory []Command

// Config : Base class to access root configuration
type Config struct {
	ID       string           `json:"id"`       // Unique ID
	Location string           `json:"location"` // Location name that refers to photo location
	Compress CompressCategory `json:"compress"` // Compression commands
	Backup   BackupCategory   `json:"backup"`   // Backup commands
}

// NewConfig build barebones data to get started on a new config file
func NewConfig(location string) *Config {
	newConfig := new(Config)               // Create empty config, and add some default info to assist in fleshing out properly
	newConfig.ID = xid.New().String()      // Generate random ID
	newConfig.Location = location          // Nice name for location
	newConfig.Compress = CompressCategory{ // Useful default entry to demo
		Command{"*.jpg *.jpeg *.png", `echo "command to run on image!"`}}
	newConfig.Backup = BackupCategory{ // Another useful demo
		Command{"harddrive", `echo "command to backup to 'harddrive'"`}}
	return newConfig
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

// Save : Save config data out for writing
func (conf *Config) Save(writer io.Writer) error {
	data, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

// Command functionality

// GetName : Get the name associated with the command
func (command Command) GetName() string {
	return command[0]
}

// GetCommand : Get the raw command string associated with the command
func (command Command) GetCommand() string {
	return command[1]
}

// Compress functionality

// GetCommand : Get the first command (in config order) whose name filter satisfies filename
func (compress CompressCategory) GetCommand(filename string) string {
	lowName := filepath.ToSlash(strings.ToLower(filename))
	for _, command := range compress {
		for _, pattern := range strings.Split(command.GetName(), " ") {
			match, err := path.Match(pattern, lowName)
			if err != nil { // This will only trigger if filter is malformed, so we should exit
				panic(err)
			}
			if match {
				return command.GetCommand()
			}
		}
	}
	return ""
}

// Backup functionality

// GetCommands : Get all backup commands that match the provided name
func (backup BackupCategory) GetCommands(name string) []string {
	commands := []string{}
	for _, command := range backup {
		match, err := filepath.Match(name, command.GetName())
		if err != nil {
			panic(err) // Malformed name!
		}
		if match {
			commands = append(commands, command.GetCommand())
		}
	}
	return commands
}
