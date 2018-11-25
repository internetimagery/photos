package rename

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func TestRename(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "18-02-01 event")
	tu.ModTime(2018, 02, 10, filepath.Join(event, "someimage.img"))

	// Get context
	cxt := tu.MustFatal(context.NewContext(event)).(*context.Context)

	// Add copy command to windows
	if runtime.GOOS == "windows" {
		cxt.Env["PATH"] = tu.Dir + ";" + cxt.Env["PATH"]
	}

	// Perform rename with compression
	tu.Must(Rename(cxt, true))

	// expecting these files
	tu.AssertExists(
		filepath.Join(event, "18-02-01 event_001.jpg"),
		filepath.Join(event, "18-02-01 event_002[one two].jpg"),
		filepath.Join(event, "18-02-10 event_003.img"),
	)
}

func TestRenameNoNew(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")

	// Get context
	cxt := tu.MustFatal(context.NewContext(event)).(*context.Context)

	// Perform rename with compression
	tu.Must(Rename(cxt, true))

	// expecting these files
	tu.AssertExists(
		filepath.Join(event, "18-12-12 event01_001.img"),
		filepath.Join(event, "18-10-11 event01_003[tags].img"),
	)
}

func TestRenameCompressCheck(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	cxt := tu.MustFatal(context.NewContext(tu.Dir)).(*context.Context)
	cxt.Env["MOCKPATH"] = filepath.Join(cxt.Root, "mockimg.jpg")
	// Add copy command to windows
	if runtime.GOOS == "windows" {
		cxt.Env["PATH"] = tu.Dir + ";" + cxt.Env["PATH"]
	}

	// Test compressed but same image
	sameimg := filepath.Join(tu.Dir, "sameimg")
	cxt.WorkingDir = sameimg

	// Perform rename with compression
	tu.Must(Rename(cxt, true))
	tu.AssertExists(filepath.Join(sameimg, "18-01-06 sameimg_001.jpg"))

	// Test compressed but same image
	diffimg := filepath.Join(tu.Dir, "diffimg")
	cxt.WorkingDir = diffimg

	// Perform rename with compression
	if err := Rename(cxt, true); err == nil {
		tu.Fail("Allowed corrupt image from third party compression")
	}
	tu.AssertNotExists(filepath.Join(diffimg, "18-01-06 diffimg_001.jpg"))
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
