// Generate or go hunting for a configuration file.
package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
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
type Command [2]string

// CompressCategory : Groups categories together. Facilitates finding commands by filter
type CompressCategory []Command

// BackupCategory : Groups backups together. Facilitates finding commands by name
type BackupCategory []Command

// Config : Base class to access root configuration
type Config struct {
	Compress CompressCategory `json:"compress"` // Compression commands
	Backup   BackupCategory   `json:"backup"`   // Backup commands
}

// NewConfig build barebones data to get started on a new config file
func NewConfig(writer io.Writer) error {
	newConfig := new(Config) // Create empty config, and add some default info to assist in fleshing out properly
	newConfig.Compress = CompressCategory{
		Command{"*.example2 *.example2", "// command to run on files ending with '.example1' or '.example2'"}}
	newConfig.Backup = BackupCategory{
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
			match, err := filepath.Match(pattern, lowName)
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
