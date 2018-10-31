package rename

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func TestRename(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Get context
	event := filepath.Join(tu.Dir, "18-02-01 event")
	cxt, err := context.NewContext(event)
	if err != nil {
		tu.Fatal(err)
	}

	// Add copy command to windows
	if runtime.GOOS == "windows" {
		cxt.Env["PATH"] = tu.Dir + ";" + cxt.Env["PATH"]
	}

	// Perform rename with compression
	if err := Rename(cxt, true); err != nil {
		tu.FailNow(err)
	}

	// expecting these files
	tu.AssertExists(
		filepath.Join(event, "18-02-01 event_001.jpg"),
		filepath.Join(event, "18-02-01 event_002[one two].jpg"),
		filepath.Join(event, "18-02-01 event_003.img"),
	)
}

func TestRenameNoNew(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Get context
	event := filepath.Join(tu.Dir, "event01")
	cxt, err := context.NewContext(event)
	if err != nil {
		fmt.Println("HERE")
		tu.Fatal(err)
	}

	// Perform rename with compression
	if err := Rename(cxt, true); err != nil {
		tu.FailNow(err)
	}

	// expecting these files
	tu.AssertExists(
		filepath.Join(event, "event01_001.img"),
		filepath.Join(event, "event01_003[tags].img"),
	)
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
