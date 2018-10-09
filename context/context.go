package context

import (
	"os"
	"path/filepath"

	"github.com/internetimagery/photos/config"
)

// ROOTCONF : name of config file that marks the root of the project (as well as important information)
const ROOTCONF = "photos.config"

// Context : Collect and encapsulate information about project
type Context struct {
	Root       string         // Path to base of repository (location of config)
	WorkingDir string         // Path of current working directory
	Config     *config.Config // Configuration information
}

// NewContext : Create a new context, gathering information
func NewContext(workingDir string) (*Context, error) {
	// Walk our way up till we find a config file denoting the project root
	currentRoot := workingDir
	configPath := ""
	for {
		configPath = filepath.Join(currentRoot, ROOTCONF) // Check if config file exists
		if _, err := os.Stat(configPath); err == nil {
			break
		}
		if currentRoot[len(currentRoot)-1] == filepath.Separator { // Check if we have hit root finding nothing, leave if so
			return nil, os.ErrNotExist
		}
		currentRoot = filepath.Dir(currentRoot)
	}

	// Attempt to load the config file
	handle, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	conf, err := config.LoadConfig(handle)
	if err != nil {
		return nil, err
	}
	return &Context{Root: currentRoot, WorkingDir: workingDir, Config: conf}, nil
}

// GetEnv : Prep environment variables to prepare commands
func (cxt *Context) GetEnv(sourcePath, destPath string) func(string) string {
	return func(name string) string {
		switch name {
		case "SOURCEPATH":
			return sourcePath
		case "DESTPATH":
			return destPath
		case "WORKINGPATH":
			return cxt.WorkingDir
		case "PROJECTPATH":
			return cxt.Root
		}
		return ""
	}
}
