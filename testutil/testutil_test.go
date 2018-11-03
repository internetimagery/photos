package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestMock : Fake tester
type TestMock struct {
	NameStr string
	LogStr  string
	Failed  bool
}

func (mock *TestMock) Name() string {
	return mock.NameStr
}
func (mock *TestMock) Fatal(message ...interface{}) {
	mock.LogStr = fmt.Sprint(message...)
	mock.Failed = true
}
func (mock *TestMock) Fail() {
	mock.Failed = true
}
func (mock *TestMock) FailNow() {
	mock.Failed = true
}
func (mock *TestMock) Log(message ...interface{}) {
	mock.LogStr = fmt.Sprint(message...)
}

func TestFail(t *testing.T) {
	mock := &TestMock{}
	tu := NewTestUtil(mock)
	tu.Fail("ARGH")
	if mock.LogStr != "ARGH" {
		t.Log("Failed to output log")
		t.Fail()
	}
}

func TestMust(t *testing.T) {
	mock := &TestMock{}
	tu := NewTestUtil(mock)
	tu.Must(nil)
	if mock.Failed {
		t.Log("Unexpected fail!")
		t.Fail()
	}
	mock.Failed = false
	tu.Must(fmt.Errorf("Fail"))
	if !mock.Failed {
		t.Log("Did not fail when supposed to")
		t.Fail()
	}
	if mock.LogStr != "Fail" {
		t.Log("Did not output message")
		t.Fail()
	}
}

// LoadTestdata : Load in testdata for testing
func TestLoadTestdata(t *testing.T) {
	tu := NewTestUtil(t)

	// Create test testdata
	close := tu.LoadTestdata()

	// Check file exists in testdata
	testFile := filepath.Join(tu.Dir, "test.file")
	if _, err := os.Stat(testFile); err != nil {
		t.Log(err)
		t.Fail()
	}

	// Cleanup testdata and check it's gone
	close()
	if _, err := os.Stat(tu.Dir); !os.IsNotExist(err) {
		if err == nil {
			t.Log("Tempfile was not removed")
		} else {
			t.Log(err)
		}
		t.Fail()
	}
}

func TestLoadTestdataMissing(t *testing.T) {
	mock := &TestMock{NameStr: "Not here!"}
	tu := NewTestUtil(mock)
	defer tu.LoadTestdata()()

	if !mock.Failed {
		t.Log("Test failed to... well... fail!")
		t.Fail()
	}
}

func TestExists(t *testing.T) {
	mock := &TestMock{NameStr: "TestExists"}
	tu := NewTestUtil(mock)
	defer tu.LoadTestdata()()

	testFile1 := filepath.Join(tu.Dir, "test.file")
	testFile2 := filepath.Join(tu.Dir, "not-test.file")

	// Filepath that exists
	tu.AssertExists(testFile1)
	if mock.Failed {
		t.Log("Test failed on existing file!")
		t.Fail()
	}
	mock.Failed = false // Reset
	tu.AssertNotExists(testFile1)
	if !mock.Failed {
		t.Log("Test succeeded on existing file!")
		t.Fail()
	}

	mock.Failed = false // Reset
	tu.AssertExists(testFile2)
	if !mock.Failed {
		t.Log("Test succeeded on missing file...")
		t.Fail()
	}
	mock.Failed = false // Reset
	tu.AssertNotExists(testFile2)
	if mock.Failed {
		t.Log("Test failed on missing file...")
		t.Fail()
	}

}

func TestModTime(t *testing.T) {
	tu := NewTestUtil(t)
	defer tu.LoadTestdata()()

	testFile := filepath.Join(tu.Dir, "test.file")
	tu.ModTime(2010, 10, 10, testFile)

	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatal(err)
	}
	testDate := info.ModTime()

	if testDate.Year() != 2010 || testDate.Month() != 10 || testDate.Day() != 10 {
		expect := time.Date(2010, 10, 10, 0, 0, 0, 0, time.UTC)
		t.Log("Expected", expect)
		t.Log("Got", testDate)
		t.Fail()
	}
}

func TestUserInput(t *testing.T) {
	tu := NewTestUtil(t)
	testMessage := "Hello"
	defer tu.UserInput(testMessage + "\n")()

	resultMessage := ""
	if _, err := fmt.Scanln(&resultMessage); err != nil {
		t.Log(err)
		t.Fail()
	}

	if resultMessage != testMessage {
		t.Log("Expected", testMessage)
		t.Log("Got", resultMessage)
		t.Fail()
	}

}
