// Sandbox utility for testing. Copy assets over to temporary test folder

package sandbox

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// SandBox : Temporary location, where tests can make or break things however they see fit
type SandBox struct {
	Root string     // Base temporary file housing the media
	t    *testing.T // Reference to the test, so we can fail if things break without having to return errors
}

// NewSandBox : Generate a new clean area to mess around with. Remember to defer "Close" to clean up.
func NewSandBox(t *testing.T) *SandBox {
	// Get location of this path, so we can get to our assets.
	_, root, _, _ := runtime.Caller(0)
	assets := filepath.Join(filepath.Dir(root), "assets")

	// Create temp dir
	tmpDir, err := ioutil.TempDir("", "Photos")
	if err != nil {
		t.Fatal(err)
	}

	// Copy files over, using "cp" command for simplicity.
	command := exec.Command("cp", "-avT", assets, tmpDir+string(filepath.Separator))
	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		t.Fatal(err)
	}

	return &SandBox{Root: tmpDir, t: t}
}

// Close : Clean up sandbox when done. Use with defer after having created in tests.
func (sb *SandBox) Close() {
	err := os.RemoveAll(sb.Root)
	if err != nil {
		sb.t.Fatal(err)
	}
}

// Join : Simple utility for joining paths to the sandbox source!
func (sb *SandBox) Join(path string) string {
	return filepath.Join(sb.Root, path)
}
