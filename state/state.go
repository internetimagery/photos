package state

import (
  "fmt"
  "github.com/internetimagery/photos/state"
  "github.com/internetimagery/photos/config"
)

type State struct {
  Root, Cwd string
  Conf *config.Config
}

func (_ State) New(cwd string) *State {
  res := new(State)
  res.Conf = new(conf.Config)
  res.Cwd = cwd
}
