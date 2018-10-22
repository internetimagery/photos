package backup

import (
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func TestRunBackup(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")
	cxt, err := context.NewContext(event)
	if err != nil {
		tu.Fatal(err)
	}

	// Set environment
	cxt.Env["TESTPATH1"] = filepath.Join(event, "testfile1.txt")
	cxt.Env["TESTPATH2"] = filepath.Join(event, "testfile2.txt")

	// Test missing command
	if err = RunBackup(cxt, "nocommand"); err != nil {
		tu.Fail(err)
	}

	// Test no command
	if err = RunBackup(cxt, ""); err != nil {
		tu.Fail(err)
	}

	// Test command
	if err = RunBackup(cxt, "test"); err != nil {
		tu.Fail(err)
	}

	// File should now exist
	tu.AssertExists(cxt.Env["TESTFILE1"])

	// Test command star
	if err = RunBackup(cxt, "othe*"); err != nil {
		tu.Fail(err)
	}

	// File should now exist
	tu.AssertExists(cxt.Env["TESTFILE2"])

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
