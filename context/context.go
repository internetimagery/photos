package context

import (
	"log"
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

// GetEnv : Expand variables in commands
func (cxt *Context) GetEnv(sourcePath, destPath string) func(string) string {
	env := map[string]string{
		"SOURCEPATH":  sourcePath,
		"DESTPATH":    destPath,
		"ROOTPATH":    cxt.Root,
		"WORKINGPATH": cxt.WorkingDir,
	}
	return func(name string) string {
		return strings.Replace(env[name], `\`, `\\`, -1)
	}
}

// RunCommand : Helper to run commands, linking outputs to terminal outputs and replacing variables safely
func RunCommand(commandString string) error {
	commandParts, err := shlex.Split(commandString)
	if err != nil {
		return err
	}
	command := exec.Command(commandParts[0], commandParts[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	log.Println("Running:", commandParts)
	return command.Run()
}
