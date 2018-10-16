package backup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/config"
	"github.com/internetimagery/photos/context"
)

func TestRunBackup(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestRunBackup")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile1 := filepath.Join(tmpDir, "testfile1.txt")
	testFile2 := filepath.Join(tmpDir, "testfile2.txt")

	cxt := &context.Context{
		Env: map[string]string{
			"TESTPATH1": testFile1,
			"TESTPATH2": testFile2,
		},
		Root:       tmpDir,
		WorkingDir: tmpDir,
		Config: &config.Config{
			Backup: config.BackupCategory{
				config.Command{"test", "touch $TESTPATH1"},
				config.Command{"other", "touch $TESTPATH2"},
			},
		},
	}

	// Test missing command
	err = RunBackup(cxt, "nocommand")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Test no command
	err = RunBackup(cxt, "")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Test command
	err = RunBackup(cxt, "test")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// File should now exist
	if _, err = os.Stat(testFile1); err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Test command star
	err = RunBackup(cxt, "othe*")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// File should now exist
	if _, err = os.Stat(testFile2); err != nil {
		fmt.Println(err)
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
			fmt.Println("Expected", value, "Got", cxt.Env[name], "from key", name)
			t.Fail()
		}
	}

}
