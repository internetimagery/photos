package config

import (
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
    "fmt"

	"github.com/rs/xid"
	"gopkg.in/yaml.v2"
)

// Command : Structure for a command
type Command struct {
	Name string `yaml:"name"`
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
    Unsorted string           `yaml:"unsorted`  // Location of folder that contains unsorted media (initial place to put media)
	Compress CompressCategory `yaml:"compress"` // Compression commands
	Backup   BackupCategory   `yaml:"backup"`   // Backup commands
}

// NewConfig build barebones data to get started on a new config file
func NewConfig(location string) *Config {
	newConfig := new(Config)               // Create empty config, and add some default info to assist in fleshing out properly
	newConfig.ID = xid.New().String()      // Generate random ID
	newConfig.Location = location          // Nice name for location
    newConfig.Unsorted = "Unsorted"        // Default location for new media
	newConfig.Compress = CompressCategory{ // Useful default entry to demo structure
		Command{Name:"*.jpg *.jpeg *.png", Command:`echo "command to run on image!"`}}
	newConfig.Backup = BackupCategory{ // Another useful demo
		Command{Name:"harddrive", Command:`echo "command to backup to 'harddrive'"`}}
	return newConfig
}

// ValidateConfig : Run some basic validations on the data
func (conf *Config) ValidateConfig() error {
    if strings.TrimSpace(conf.Location) == "" {
        return fmt.Errorf("empty project Location Name")
    }
    trimUnsorted := path.Clean(strings.TrimSpace(conf.Unsorted))
    if trimUnsorted == "" {
        return fmt.Errorf("unsorted path is empty")
    }
    if path.IsAbs(trimUnsorted) {
        return fmt.Errorf("unsorted path is absolute, must be relative")
    }
    if trimUnsorted == "." || strings.HasPrefix(trimUnsorted, "..")  {
        return fmt.Errorf("unsorted path must be within project")
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
		return err
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
