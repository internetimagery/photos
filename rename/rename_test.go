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
	// Mock context
	mockCxt := &context.Context{
		Config: &config.Config{
			Compress: config.CompressCategory{
				config.Command{"*", `echo "$SOURCEPATH"`},
			},
		},
	}

	// Working Path
	tmpDir, err := ioutil.TempDir("", "photo_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	eventName := "18-02-01 event"
	rootPath := filepath.Join(tmpDir, eventName)
	err = os.Mkdir(rootPath, 755)
	if err != nil {
		t.Fatal(err)
	}

	// Prep some test files
	testFiles := map[string]string{
		"someimage.img":                 eventName + "_003.img",
		eventName + "_001.jpg":          eventName + "_001.jpg",
		eventName + "_002[one two].jpg": eventName + "_002[one two].jpg",
	}

	// Create files
	for name := range testFiles {
		err = ioutil.WriteFile(filepath.Join(rootPath, name), []byte("some data"), 655)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Perform rename
	err = Rename(rootPath, mockCxt)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Check files made it to where they need to be

	// Check source files are where they should be
	sourcePath := filepath.Join(tmpDir, SOURCEDIR)
	for src, dst := range testFiles {
		if src != dst {
			if _, err = os.Stat(filepath.Join(sourcePath, src)); err != nil {
				fmt.Println(err)
				t.Fail()
			}
		}
	}
}
