// Drop command
package commands

import (
  "log"
  "github.com/internetimagery/photos/state"
)

type CMD_Drop struct{}

func (self CMD_Drop) Desc() string {
  return "Drop content, replacing with placeholder."
}

func (self CMD_Drop) Run(args []string, state *state.State) int {
  log.Println(args)
  log.Println("DROP")
  log.Println(state)
  return 0
}
