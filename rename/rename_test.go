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

	// Mock context
	mockCxt := &context.Context{
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
		err = ioutil.WriteFile(filepath.Join(rootPath, name), []byte("some data"), 655)
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
