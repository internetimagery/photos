package backup

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func TestRunBackup(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")
	cxt := tu.MustFatal(context.NewContext(event)).(*context.Context)

	// Add touch command to windows
	if runtime.GOOS == "windows" {
		cxt.Env["PATH"] = tu.Dir + ";" + cxt.Env["PATH"]
	}

	// Set environment
	today := time.Now().Format("06-01-02")
	cxt.Env["TESTPATH1"] = filepath.Join(event, today+" event01_001.txt")
	cxt.Env["TESTPATH2"] = filepath.Join(event, today+" event01_002.txt")
	cxt.Env["TESTPATH3"] = filepath.Join(event, today+" event01_003.txt")

	// Test missing command
	tu.Must(RunBackup(cxt, "nocommand"))

	// Test no command
	tu.Must(RunBackup(cxt, ""))

	// Test command
	tu.Must(RunBackup(cxt, "test"))

	// File should now exist
	tu.AssertExists(cxt.Env["TESTPATH1"])

	// Test command star
	tu.Must(RunBackup(cxt, "othe*"))

	// Files should now exist
	tu.AssertExists(cxt.Env["TESTPATH2"])
	tu.AssertExists(cxt.Env["TESTPATH3"])

}

func TestBackupBadCommand(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")
	cxt := tu.MustFatal(context.NewContext(event)).(*context.Context)

	if err := RunBackup(cxt, "test"); err == nil {
		tu.Fail("Passed on bad command!")
	}
}

func TestBackupContainsSource(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")
	cxt := tu.MustFatal(context.NewContext(event)).(*context.Context)

	if err := RunBackup(cxt, "test"); err == nil {
		tu.Fail("Allowed backup with source files still present.")
	}
}

func TestBackupContainsUnformatted(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")
	cxt := tu.MustFatal(context.NewContext(event)).(*context.Context)

	if err := RunBackup(cxt, "test"); err == nil {
		tu.Fail("Allowed backup with source files still unformatted.")
	}
}

func TestSetEnviron(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	working := "/path/to/files"
	root := "/path"
	relworking := "to/files"

	cxt := &context.Context{
		Env:        map[string]string{},
		WorkingDir: working,
		Root:       root,
	}

	// Set up our environment
	setEnvironment(cxt)

	testCase := map[string]string{
		"SOURCEPATH":  working,
		"ROOTPATH":    root,
		"WORKINGPATH": working,
		"RELPATH":     relworking,
	}

	for name, value := range testCase {
		if cxt.Env[name] != value {
			tu.FailE(value, cxt.Env[name])
		}
	}
}
