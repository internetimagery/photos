package backup

import (
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/config"
	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func TestRunBackup(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.TempDir("TestRunBackup")()

	testFile1 := filepath.Join(tu.Dir, "testfile1.txt")
	testFile2 := filepath.Join(tu.Dir, "testfile2.txt")

	cxt := &context.Context{
		Env: map[string]string{
			"TESTPATH1": testFile1,
			"TESTPATH2": testFile2,
		},
		Root:       tu.Dir,
		WorkingDir: tu.Dir,
		Config: &config.Config{
			Backup: config.BackupCategory{
				config.Command{"test", "touch $TESTPATH1"},
				config.Command{"other", "touch $TESTPATH2"},
			},
		},
	}

	// Test missing command
	err := RunBackup(cxt, "nocommand")
	if err != nil {
		tu.Fail(err)
	}

	// Test no command
	err = RunBackup(cxt, "")
	if err != nil {
		tu.Fail(err)
	}

	// Test command
	err = RunBackup(cxt, "test")
	if err != nil {
		tu.Fail(err)
	}

	// File should now exist
	tu.AssertExists(testFile1)

	// Test command star
	err = RunBackup(cxt, "othe*")
	if err != nil {
		tu.Fail(err)
	}

	// File should now exist
	tu.AssertExists(testFile2)

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
