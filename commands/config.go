// Config command
package commands

import (
  "log"
  "github.com/internetimagery/photos/state"
)

type CMD_Config struct{}

func (self CMD_Config) Desc() string {
  return "Manage configuration."
}

func (self CMD_Config) Run(args []string, state *state.State) int {
  log.Println(args)
  log.Println("CONFIG")
  log.Println(state)
  return 0
}
