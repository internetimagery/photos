package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// LoadTestdata : Load in testdata for testing
func TestLoadTestdata(t *testing.T) {
	tu := NewTestUtil(t)

	// Create test testdata
	close := tu.LoadTestdata()

	// Check file exists in testdata
	testFile := filepath.Join(tu.Dir, "test.file")
	if _, err := os.Stat(testFile); err != nil {
		tu.Fail(err)
	}

	// Cleanup testdata and check it's gone
	close()
	if _, err := os.Stat(tu.Dir); !os.IsNotExist(err) {
		if err == nil {
			tu.Fail("Tempfile was not removed")
		} else {
			tu.Fail(err)
		}
	}
}

func TestExists(t *testing.T) {
	tu := NewTestUtil(t)
	defer tu.LoadTestdata()()

	testFile1 := filepath.Join(tu.Dir, "test.file")
	testFile2 := filepath.Join(tu.Dir, "not-test.file")
	tu.AssertExists(testFile1)
	tu.AssertNotExists(testFile2)
}

func TestNew(t *testing.T) {
	tu := NewTestUtil(t)
	defer tu.LoadTestdata()()

	testFile := filepath.Join(tu.Dir, "test.file")
	testDir := filepath.Join(tu.Dir, "test-dir")

	tu.NewFile(testFile, "testing")
	tu.NewDir(testDir)

	tu.AssertExists(testFile)
	tu.AssertExists(testDir)
}

func TestTempDir(t *testing.T) {
	tu := NewTestUtil(t)
	close := tu.TempDir("TestTempDir")
	if _, err := os.Stat(tu.Dir); err != nil {
		tu.Fail(err)
	}
	close()
	if _, err := os.Stat(tu.Dir); !os.IsNotExist(err) {
		tu.Fail("Tempdir not removed.")
	}

}

func TestUserInput(t *testing.T) {
	tu := NewTestUtil(t)
	testMessage := "Hello"
	defer tu.UserInput(testMessage + "\n")()

	resultMessage := ""
	if _, err := fmt.Scanln(&resultMessage); err != nil {
		tu.Fail(err)
	}

	if resultMessage != testMessage {
		tu.FailE(testMessage, resultMessage)
	}

}
