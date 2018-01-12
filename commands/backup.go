// Backup command
package commands

import (
  "log"
  "github.com/internetimagery/photos/config"
)

type CMD_Backup struct{}

func (self CMD_Backup) Desc() string {
  return "Copy data to a remote location."
}

func (self CMD_Backup) Run(args []string, conf *config.Config) int {
  log.Println(args)
  log.Println("BACKUP")
  log.Println(conf)
  return 0
}
