// Initialize command
package commands

import (
  "log"
  "github.com/internetimagery/photos/config"
)

type CMD_Init struct{}

func (self CMD_Init) Desc() string {
  return "Initialize a new repository."
}

func (self CMD_Init) Run(args []string, conf *config.Config) int {
  log.Println(args)
  log.Println("INIT")
  log.Println(conf)
  return 0
}
