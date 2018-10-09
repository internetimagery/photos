package context

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
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
	configData := []byte(`{
	"compress": [
		["*", "some command"]
	]
}`)

	tmpDir, err := ioutil.TempDir("", "photos-context-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Make a couple files
	workingDir := filepath.Join(tmpDir, "some-event")
	rootConf := filepath.Join(tmpDir, ROOTCONF)

	err = os.Mkdir(workingDir, 755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(rootConf, configData, 655)
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

func TestGetEnv(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "photo_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err = ioutil.WriteFile(filepath.Join(tmpDir, ROOTCONF), []byte("{}"), 644); err != nil {
		t.Fatal(err)
	}

	cxt, err := NewContext(tmpDir)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	sourcePath := filepath.Join(tmpDir, "source")
	destPath := filepath.Join(tmpDir, "dest")
	envList := cxt.GetEnv(sourcePath, destPath)

	testEnv := map[string]bool{
		"SOURCEPATH=" + sourcePath: true,
		"DESTPATH=" + destPath:     true,
		"WORKINGPATH=" + tmpDir:    true,
		"PROJECTPATH=" + tmpDir:    true,
	}

	for _, env := range envList {
		if !testEnv[env] {
			fmt.Println("Environment incorrect", env)
			t.Fail()
		}
	}
	if len(envList) != len(testEnv) {
		fmt.Println("Envlist does not match expected")
		fmt.Println("Expected:")
		fmt.Println(testEnv)
		fmt.Println("Got:")
		fmt.Println(testEnv)
		t.Fail()
	}

}
