package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFile(t *testing.T) {
	tu := NewTestUtil(t)
	defer tu.TempDir("TestNewFile")()

	testFile := filepath.Join(tu.Dir, "test.file")
	tu.NewFile(testFile, "")
	if _, err := os.Stat(testFile); err != nil {
		tu.Fail(err)
	}
}

func TestExists(t *testing.T) {
	tu := NewTestUtil(t)
	defer tu.TempDir("TestExists")()

	testFile := filepath.Join(tu.Dir, "test.file")
	tu.NewFile(testFile, "")
	tu.AssertExists(testFile)
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
