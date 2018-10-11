package context

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewContext(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestNewContext")
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
	configData := []byte(`{
	"compress": [
		["*", "some command"]
	]
}`)

	tmpDir, err := ioutil.TempDir("", "TestContext")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Make a couple files
	workingDir := filepath.Join(tmpDir, "some-event")
	rootConf := filepath.Join(tmpDir, ROOTCONF)

	err = os.Mkdir(workingDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(rootConf, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Start within event directory
	cxt, err := NewContext(workingDir)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Check we reached root file
	if cxt.Root != tmpDir {
		fmt.Println("Couldn't find config file.", cxt.Root)
		t.Fail()
	}
}
