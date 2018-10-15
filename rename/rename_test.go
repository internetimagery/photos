package rename

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/config"
	"github.com/internetimagery/photos/context"
)

func TestRename(t *testing.T) {

	// Working Path
	tmpDir, err := ioutil.TempDir("", "TestRename")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	eventName := "18-02-01 event"
	rootPath := filepath.Join(tmpDir, eventName)
	err = os.Mkdir(rootPath, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Mock context
	mockCxt := &context.Context{
		Env:        map[string]string{},
		WorkingDir: rootPath,
		Config: &config.Config{
			Compress: config.CompressCategory{
				config.Command{"*", `cp "$SOURCEPATH" "$DESTPATH"`},
			},
		},
	}

	// Prep some test files
	testFiles := map[string]string{
		"someimage.img":                 eventName + "_003.img",
		eventName + "_001.jpg":          eventName + "_001.jpg",
		eventName + "_002[one two].jpg": eventName + "_002[one two].jpg",
	}

	// Create files
	for name := range testFiles {
		err = ioutil.WriteFile(filepath.Join(rootPath, name), []byte("some data"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Perform rename with compression
	err = Rename(mockCxt, true)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Check files made it to where they need to be
	sourcePath := filepath.Join(rootPath, SOURCEDIR)
	for src, dst := range testFiles {
		if src != dst {
			// Check renamed
			if _, err = os.Stat(filepath.Join(rootPath, dst)); err != nil {
				fmt.Println(err)
				t.Fail()
			}
			// Check original source
			if _, err = os.Stat(filepath.Join(sourcePath, src)); err != nil {
				fmt.Println(err)
				t.Fail()
			}
		}
	}
}

func TestSetEnviron(t *testing.T) {
	cxt := &context.Context{Env: map[string]string{}}

	sourcePath := "/path/to/original.file"
	destPath := "/path/to/other.file"

	// Set up our environment
	setEnviron(sourcePath, destPath, cxt)

	testCase := map[string]string{
		"SOURCEPATH":  sourcePath,
		"DESTPATH":    destPath,
		"ROOTPATH":    cxt.Root,
		"WORKINGPATH": cxt.WorkingDir,
	}

	for name, value := range testCase {
		if cxt.Env[name] != value {
			fmt.Println("Expected", value, "Got", cxt.Env[name], "from key", name)
			t.Fail()
		}
	}

}
