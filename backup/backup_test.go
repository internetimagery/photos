package backup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
		Root:       tmpDir,
		WorkingDir: tmpDir,
		Config: &config.Config{
			Backup: config.BackupCategory{
				config.Command{"test", fmt.Sprintf("touch '%s'", strings.Replace(testFile1, `\`, `\\`, -1))},
				config.Command{"other", fmt.Sprintf("touch '%s'", strings.Replace(testFile2, `\`, `\\`, -1))},
			},
		},
	}

	// Test missing command
	err = RunBackup(cxt, "nocommand")
	if err == nil {
		fmt.Println("Missing command returned no error", "nocommand")
		t.Fail()
	}

	// Test no command
	err = RunBackup(cxt, "")
	if err == nil {
		fmt.Println("Empty command returned no error")
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
