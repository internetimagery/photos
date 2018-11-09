package context

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
	"github.com/internetimagery/photos/config"
)

// ROOTCONF : name of config file that marks the root of the project (as well as important information)
const ROOTCONF = "photos-config.json"

// Context : Collect and encapsulate information about project
type Context struct {
	Root       string            // Path to base of repository (location of config)
	WorkingDir string            // Path of current working directory
	Env        map[string]string // Representation of the environment
	Config     *config.Config    // Configuration information
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

	// Build out our environment vars
	env := make(map[string]string)
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if parts[0] != "" {
			env[parts[0]] = parts[1]
		}
	}

	return &Context{Root: currentRoot, WorkingDir: workingDir, Config: conf, Env: env}, nil
}

// expandEnv : Expand environment variables with those from context. Make safe the backslashes also!
func (cxt *Context) expandEnv(name string) string {
	return strings.Replace(cxt.Env[name], `\`, `\\`, -1)
}

// contractEnv : Turn environment map back into "KEY=VALUE" pairs for use in commands
func (cxt *Context) contractEnv() []string {
	env := []string{}
	for key, value := range cxt.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}

// PrepCommand : Take string. Prepare a command object from it, and expand variables with their environment counterparts
func (cxt *Context) PrepCommand(commandRaw string) (*exec.Cmd, error) {
	commandExpanded := os.Expand(commandRaw, cxt.expandEnv) // Expand variables in string name
	commandParts, err := shlex.Split(commandExpanded)       // Split command into parts for execution
	if err != nil {
		return nil, err
	}
	command := exec.Command(commandParts[0], commandParts[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = cxt.contractEnv()
	return command, nil
}


// Done