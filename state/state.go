package state

import (
	"log"
	"os"
	"path/filepath"

	"github.com/internetimagery/photos/config"
	"github.com/internetimagery/photos/utility"
)

const CONFIG = "test_config.json"

// Manage basic cli state
type State struct {
	Root, Cwd, Conf_Path string
	Conf                 *config.Config
}

func (_ State) New() *State {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	// Look for config file
	conf := utility.SearchUp(CONFIG, cwd)
	root := filepath.Dir(conf)

	// TODO: Collect config path and file etc
	state := new(State)
	state.Cwd = cwd
	state.Conf_Path = conf
	state.Root = root
	return state
}
