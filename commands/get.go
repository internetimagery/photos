// Get command
package commands

import (
  "log"
  "github.com/internetimagery/photos/state"
)

type CMD_Get struct{}

func (self CMD_Get) Desc() string {
  return "Get placeholder data."
}

func (self CMD_Get) Run(args []string, state *state.State) int {
  log.Println(args)
  log.Println("GET")
  log.Println(state)
  return 0
}
