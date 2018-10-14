package backup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/internetimagery/photos/context"
)

func TestRunBackup(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestRunBackup")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "testfile.txt")

	configData := []byte(fmt.Sprintf(`{
    "backup":
      ["test", "touch \"%s\""]
    }`, strings.Replace(testFile, `\`, `\\`, -1)))
	err = ioutil.WriteFile(filepath.Join(tmpDir, context.ROOTCONF), configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	cxt, err := context.NewContext(tmpDir)
	if err != nil {
		fmt.Println(string(configData))
		fmt.Println(err)
		t.Fail()
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
	if _, err = os.Stat(testFile); err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
