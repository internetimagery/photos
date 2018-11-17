package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/xid"
	"gopkg.in/yaml.v2"
)

// SORTED : Default path to file where sorted media goes (before being assigned an event or being renamed/compressed)
const SORTED = "Sorted"

// Command : Structure for a command
type Command struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

// CompressCategory : Groups categories together. Facilitates finding commands by filter
type CompressCategory []Command

// BackupCategory : Groups backups together. Facilitates finding commands by name
type BackupCategory []Command

// Config : Base class to access root configuration
type Config struct {
	ID       string           `yaml:"id"`       // Unique ID
	Location string           `yaml:"location"` // Location name that refers to project
	Sorted   string           `yaml:"sorted"`   // Location of folder that contains sorted media (before being assigned an event/compressed)
	Compress CompressCategory `yaml:"compress"` // Compression commands
	Backup   BackupCategory   `yaml:"backup"`   // Backup commands
}

// NewConfig build barebones data to get started on a new config file
func NewConfig(location string) *Config {
	newConfig := new(Config)               // Create empty config, and add some default info to assist in fleshing out properly
	newConfig.ID = xid.New().String()      // Generate random ID
	newConfig.Location = location          // Nice name for location
	newConfig.Sorted = SORTED              // Default location for sorted media
	newConfig.Compress = CompressCategory{ // Useful default entry to demo structure
		Command{Name: "*.jpg *.jpeg *.png", Command: `echo "command to run on image!"`}}
	newConfig.Backup = BackupCategory{ // Another useful demo
		Command{Name: "harddrive", Command: `echo "command to backup to 'harddrive'"`}}
	return newConfig
}

// validatePath : Helper to validate a path exists, is relative, and does not descend backwards
func validatePath(filename string) error {
	trimname := strings.TrimSpace(filename)
	if trimname == "" {
		return fmt.Errorf("path is empty")
	}
	cleanname := path.Clean(trimname)
	if path.IsAbs(cleanname) {
		return fmt.Errorf("path is absolute, must be relative '%s'", cleanname)
	}
	if cleanname == "." {
		return fmt.Errorf("path cannot be root directory")
	}
	if strings.HasPrefix(cleanname, "..") {
		return fmt.Errorf("path must be within project '%s'", cleanname)
	}
	return nil
}

// ValidateConfig : Run some basic validations on the data
func (conf *Config) ValidateConfig() error {
	if strings.TrimSpace(conf.Location) == "" {
		return fmt.Errorf("empty project Location Name")
	}
	if err := validatePath(conf.Sorted); err != nil {
		return err
	}
	return nil
}

// LoadConfig : Load and populate a new Config from existing config data
func LoadConfig(reader io.Reader) (*Config, error) {
	loadedData, err := ioutil.ReadAll(reader) // Load the data to process
	if err != nil {
		return nil, err
	}
	loadedConfig := new(Config)
	err = yaml.Unmarshal(loadedData, loadedConfig)
	if err != nil {
		return nil, err
	}
	if loadedConfig.Sorted == "" { // Set default
		loadedConfig.Sorted = SORTED
	}
	return loadedConfig, loadedConfig.ValidateConfig()
}

// Save : Save config data out for writing
func (conf *Config) Save(writer io.Writer) error {
	err := conf.ValidateConfig()
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(conf)
	if err != nil {
		panic(err) // This should never fail!
	}
	_, err = writer.Write(data)
	return err
}

// Compress functionality

// GetCommand : Get the first command (in config order) whose name filter satisfies filename
func (compress CompressCategory) GetCommand(filename string) string {
	lowName := filepath.Base(strings.ToLower(filename))
	for _, command := range compress {
		for _, pattern := range strings.Split(command.Name, " ") {
			match, err := path.Match(pattern, lowName)
			if err != nil { // This will only trigger if filter is malformed, so we should exit
				panic(err)
			}
			if match {
				return command.Command
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
		match, err := filepath.Match(name, command.Name)
		if err != nil {
			panic(err) // Malformed name!
		}
		if match {
			commands = append(commands, command.Command)
		}
	}
	return commands
}
