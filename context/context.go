package context

import "github.com/derekparker/delve/pkg/config"

// ROOTCONF : name of config file that marks the root of the project (as well as important information)
const ROOTCONF = "photos.config"

// Context : Collect and encapsulate information about project
type Context struct {
	Root             string        // Path to base of repository (location of config)
	WorkingDirectory string        // Path of current working directory
	Config           config.Config // Configuration information
}

// NewContext : Create a new context, gathering information
func NewContext(workingDirectory string) (*Context, error) {
	// TODO: search up the path from workingDirectory to find config file
	return nil, nil
}

// GetRelativePath : Make given path relative to project root. Error if not in project...
func GetRelativePath(path string) (string, error) {
	return "", nil
}
