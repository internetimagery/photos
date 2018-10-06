package context

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/internetimagery/photos/sandbox"
)

func TestNewContext(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "photostest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	if _, err := NewContext(tmpDir); !os.IsNotExist(err) {
		fmt.Println(err)
		t.Fail()
	}
}

func TestContext(t *testing.T) {
	sb := sandbox.NewSandBox(t)
	defer sb.Close()

	// Start within event directory
	cxt, err := NewContext(sb.Join("event01"))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Check root and workingDir are not the same!
	if cxt.Root == cxt.WorkingDir {
		fmt.Println("Root is the same as workingDir!")
		t.Fail()
	}
}
