// Get command
package commands

import (
  "log"
  "github.com/internetimagery/photos/config"
)

type CMD_Get struct{}

func (self CMD_Get) Desc() string {
  return "Get placeholder data."
}

func (self CMD_Get) Run(args []string, conf *config.Config) int {
  log.Println(args)
  log.Println("GET")
  log.Println(conf)
  return 0
}
