// Drop command
package commands

import (
  "log"
  "github.com/internetimagery/photos/config"
)

type CMD_Drop struct{}

func (self CMD_Drop) Desc() string {
  return "Drop content, replacing with placeholder."
}

func (self CMD_Drop) Run(args []string, conf *config.Config) int {
  log.Println(args)
  log.Println("DROP")
  log.Println(conf)
  return 0
}
