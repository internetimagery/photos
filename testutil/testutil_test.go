package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestModTime(t *testing.T) {
	tu := NewTestUtil(t)
	defer tu.LoadTestdata()()

	testFile := filepath.Join(tu.Dir, "test.file")
	tu.ModTime(2010, 10, 10, testFile)

	info, err := os.Stat(testFile)
	if err != nil {
		tu.Fatal(err)
	}
	testDate := info.ModTime()

	if testDate.Year() != 2010 || testDate.Month() != 10 || testDate.Day() != 10 {
		expect := time.Date(2010, 10, 10, 0, 0, 0, 0, time.UTC)
		tu.FailE(expect, testDate)
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
