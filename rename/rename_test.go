package rename

import (
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/config"
	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func TestRename(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	// Working Path
	defer tu.TempDir("TestRename")()

	eventName := "18-02-01 event"
	rootPath := filepath.Join(tu.Dir, eventName)
	tu.NewDir(rootPath)

	// Mock context
	mockCxt := &context.Context{
		Env:        map[string]string{},
		WorkingDir: rootPath,
		Config: &config.Config{
			Compress: config.CompressCategory{
				config.Command{"*", `cp -v "$SOURCEPATH" "$DESTPATH"`},
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
		tu.NewFile(filepath.Join(rootPath, name), "")
	}

	// Perform rename with compression
	if err := Rename(mockCxt, true); err != nil {
		tu.Fail(err)
	}

	// Check files made it to where they need to be
	sourcePath := filepath.Join(rootPath, SOURCEDIR)
	for src, dst := range testFiles {
		if src != dst {
			// Check renamed
			tu.AssertExists(filepath.Join(rootPath, dst))

			// Check original source
			tu.AssertExists(filepath.Join(sourcePath, src))
		}
	}
}

func TestSetEnviron(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	sourcePath := "/path/to/original.file"
	destPath := "/path/to/other.file"
	working := "/path/to"
	root := "/path"

	cxt := &context.Context{
		Env:        map[string]string{},
		WorkingDir: working,
		Root:       root,
	}

	// Set up our environment
	setEnvironment(sourcePath, destPath, cxt)

	testCase := map[string]string{
		"SOURCEPATH":  sourcePath,
		"DESTPATH":    destPath,
		"ROOTPATH":    root,
		"WORKINGPATH": working,
	}

	for name, value := range testCase {
		if cxt.Env[name] != value {
			tu.FailE(value, cxt.Env[name])
		}
	}
}
