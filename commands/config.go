// Config command
package commands

import (
  "log"
  "github.com/internetimagery/photos/config"
)

type CMD_Config struct{}

func (self CMD_Config) Desc() string {
  return "Manage configuration."
}

func (self CMD_Config) Run(args []string, conf *config.Config) int {
  log.Println(args)
  log.Println("CONFIG")
  log.Println(conf)
  return 0
}
