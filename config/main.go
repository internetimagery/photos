// Repo configuration
package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	uuid "github.com/satori/go.uuid"
)

// "fmt"

// RandomID from mac address and time
func GenerateID() string {
	return uuid.NewV4().String()
}

type Config struct {
	ID, Name, Root string
}

// Create a new config file
func NewConfig() *Config {
	// Perform initial setup here
	conf := new(Config)
	conf.ID = GenerateID()
	return conf
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := new(Config)
	err = json.Unmarshal(data, conf)
	if conf.ID == "" {
		return nil, errors.New("Missing ID")
	}
	conf.Root = path
	return conf, err
}

func (self Config) Save(path string) error {
	data, err := json.MarshalIndent(self, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 644)
	return err
}

//
// const CONFIGNAME = "photos_conf.json"
//
// // REMOTE!
// // UUID: Unique identifier (taken from repo config)
// // Name: Friendly name (from repo config)
// // Locations: List of prefixes that leads to repo
// type Remote struct {
//   UUID, Name string
//   Locations []string
// }
//
// // CONFIG!
// // UUID: Unique identifier
// // Name: Friendly Name
// // Rclone: rclone location (command)
// // Bin: (optional) path to rclone
// // Remotes: remote repos
// type Config struct {
// 	UUID, Name, Rclone, root string
//   Remotes map[string]*Remote
// }
// func (self Config) GetRoot() string {
//   return self.root
// }
//
// // Find config in the current context
// func GetConfig(root string) *Config {
//   var conf Config
//   pconf := &conf
//   path := utility.SearchUp(CONFIGNAME, root)
//   if path != "" {
//     pconf = LoadConfig(path)
//   }
//   return pconf
// }
//
// // Create new config
// func NewConfig() *Config {
// 	conf := &Config{UUID: utility.GenerateID()}
// 	return conf
// }
//
// // Save config file to disk
// func SaveConfig(conf *Config, path string)  {
//   data, err := json.MarshalIndent(conf, "", "  ")
//   if err != nil {
//     log.Panic(err)
//   }
//   err = ioutil.WriteFile(path, data, 644)
//   if err != nil {
//     log.Panic(err)
//   }
// }
//
// // Read a config from file
// func LoadConfig(path string) *Config {
//   data, err := ioutil.ReadFile(path)
//   if err != nil {
//     log.Panic(err)
//   }
//   conf := &Config{root:filepath.Dir(path)}
//   err = json.Unmarshal(data, conf)
//   if err != nil {
//     log.Panic(err)
//   }
//   return conf
// }
//
// // cwd := utility.CWD()
// // files := ioutil.ReadDir(cwd)
// // for i := 0; i < len(files); i++ {
// //   if files[i].Name() == CONFIGNAME {
// //     log.Fatal("Config file already exists.")
// //   }
// // }
