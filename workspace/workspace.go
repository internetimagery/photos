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
type Workspace struct {
	Root, Cwd, Conf_Path string
	Conf                 *config.Config
}

// Create new state
func (_ Workspace) New() *Workspace {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	// Look for config file
	conf := utility.SearchUp(CONFIG, cwd)
	root := filepath.Dir(conf)
	data := new(config.Config)

	// TODO: Collect config path and file etc
	ws := new(Workspace)
	ws.Cwd = cwd
	ws.Conf = data
	ws.Conf_Path = conf
	ws.Root = root
	return ws
}
