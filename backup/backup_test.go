package backup

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/config"
	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func TestRunBackup(t *testing.T) {
	tmpDir := testutil.NewTempDir(t, "TestRunBackup")
	defer tmpDir.Close()

	testFile1 := filepath.Join(tmpDir.Dir, "testfile1.txt")
	testFile2 := filepath.Join(tmpDir.Dir, "testfile2.txt")

	cxt := &context.Context{
		Env: map[string]string{
			"TESTPATH1": testFile1,
			"TESTPATH2": testFile2,
		},
		Root:       tmpDir.Dir,
		WorkingDir: tmpDir.Dir,
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
		t.Log(err)
		t.Fail()
	}

	// Test no command
	err = RunBackup(cxt, "")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// Test command
	err = RunBackup(cxt, "test")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// File should now exist
	if _, err = os.Stat(testFile1); err != nil {
		t.Log(err)
		t.Fail()
	}

	// Test command star
	err = RunBackup(cxt, "othe*")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// File should now exist
	if _, err = os.Stat(testFile2); err != nil {
		t.Log(err)
		t.Fail()
	}

}

func TestSetEnviron(t *testing.T) {

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
			t.Log("Expected", value, "Got", cxt.Env[name], "from key", name)
			t.Fail()
		}
	}

}
