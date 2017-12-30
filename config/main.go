// Repo configuration
package config

import (
  "log"
  "io/ioutil"
  "github.com/internetimagery/photos/utility"
)

const CONFIGNAME = "photos_conf.json"

type Config struct {
  UUID, Name string
}

// Find config in the current context
func GetConfig() string {
  cwd := utility.CWD()
  return utility.SearchUp(CONFIGNAME, cwd)
}

// Create new config
func NewConfig(name string) *Config {
  conf := &Config{UUID: utility.GenerateID(), Name: name}
  return conf
}

// Read a config
func LoadConfig(path string) *Config {
  
}


// cwd := utility.CWD()
// files := ioutil.ReadDir(cwd)
// for i := 0; i < len(files); i++ {
//   if files[i].Name() == CONFIGNAME {
//     log.Fatal("Config file already exists.")
//   }
// }
