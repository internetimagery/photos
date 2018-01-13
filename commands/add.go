// Add command
package commands

import (
  "log"
  "github.com/internetimagery/photos/state"
)

type CMD_Add struct{}

func (self CMD_Add) Desc() string {
  return "Add media to repository."
}

func (self CMD_Add) Run(args []string, state *state.State) int {
  log.Println(args)
  log.Println("ADD")
  log.Println(state)
  return 0
}
