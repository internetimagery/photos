// Initialize command
package commands

import (
  "log"
  "github.com/internetimagery/photos/state"
)

type CMD_Init struct{}

func (self CMD_Init) Desc() string {
  return "Initialize a new repository."
}

func (self CMD_Init) Run(args []string, state *state.State) int {
  log.Println(args)
  log.Println("INIT")
  log.Println(state)
  return 0
}
