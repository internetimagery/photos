package rename

import (
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func TestRename(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Get context
	cxt, err := context.NewContext(tu.Dir)
	if err != nil {
		tu.Fatal(err)
	}

	// Perform rename with compression
	if err := Rename(cxt, true); err != nil {
		tu.Fail(err)
	}

	// expecting these files
	tu.AssertExistsAll(
		filepath.Join(tu.Dir, "18-02-01 event", "18-02-01 event_001.jpg"),
		filepath.Join(tu.Dir, "18-02-01 event", "18-02-01 event_002[one two].jpg"),
		filepath.Join(tu.Dir, "18-02-01 event", "18-02-01 event_003.img"),
	)
}

func TestRenameNoNew(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Get context
	cxt, err := context.NewContext(tu.Dir)
	if err != nil {
		tu.Fatal(err)
	}

	// Perform rename with compression
	if err := Rename(cxt, true); err != nil {
		tu.Fail(err)
	}

	// expecting these files
	tu.AssertExistsAll(
		filepath.Join(tu.Dir, "event01", "event01_001.img"),
		filepath.Join(tu.Dir, "event01", "event01_003[tags].img"),
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
