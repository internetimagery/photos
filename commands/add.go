// Add command
package commands

import (
  "log"
  "github.com/internetimagery/photos/config"
)

type CMD_Add struct{}

func (self CMD_Add) Desc() string {
  return "Add media to repository."
}

func (self CMD_Add) Run(args []string, conf *config.Config) int {
  log.Println(args)
  log.Println("ADD")
  log.Println(conf)
  return 0
}
