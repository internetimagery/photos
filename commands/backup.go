// Backup command
package commands

import (
  "log"
  "github.com/internetimagery/photos/state"
)

type CMD_Backup struct{}

func (self CMD_Backup) Desc() string {
  return "Copy data to a remote location."
}

func (self CMD_Backup) Run(args []string, state *state.State) int {
  log.Println(args)
  log.Println("BACKUP")
  log.Println(state)
  return 0
}
