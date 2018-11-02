package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/internetimagery/photos/copy"
)

type Tester interface {
	Name() string
	Fatal(...interface{})
	Fail()
	FailNow()
	Log(...interface{})
}

// TestUtil : Wrapper for helper functions
type TestUtil struct {
	Dir string
	Tester
}

// NewTestUtil : Create new testutil
func NewTestUtil(t Tester) *TestUtil {
	return &TestUtil{Tester: t}
}

// LoadTestdata : Copy across testdata into temporary directory
func (util *TestUtil) LoadTestdata() func() {
	testdata, err := filepath.Abs(filepath.Join("testdata", util.Name()))
	if err != nil {
		util.Fatal()
	}
	if _, err = os.Stat(testdata); err != nil {
		if os.IsNotExist(err) {
			util.Fatal("Test data does not exist with name", util.Name())
		} else {
			util.Fatal(err)
		}
	}
	tmpDir, err := ioutil.TempDir("", util.Name())
	if err != nil {
		util.Fatal(err)
	}
	if err = copy.Tree(testdata, tmpDir); err != nil {
		util.Fatal(err)
	}
	util.Dir = tmpDir
	return func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			util.Fatal(err)
		}
		util.Dir = ""
	}
}

// ModTime : Set modification time of file
func (util *TestUtil) ModTime(year, month, day int, filePaths ...string) {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	for _, filePath := range filePaths {
		if err := os.Chtimes(filePath, date, date); err != nil {
			util.Fatal(err)
		}
	}
}

// AssertExists : Check if file exists. Fail if not
func (util *TestUtil) AssertExists(filePaths ...string) []os.FileInfo {
	result := []os.FileInfo{}
	for _, filePath := range filePaths {
		info, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				util.Fail("File does not exist:", filePath)
			} else {
				util.Fail(err)
			}
		}
		result = append(result, info)
	}
	return result
}

// AssertNotExists : Check if file is missing. Fail if it exists.
func (util *TestUtil) AssertNotExists(filePaths ...string) {
	for _, filePath := range filePaths {
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			if err == nil {
				util.Fail("File exists, but shouldn't:", filePath)
			} else {
				util.Fail(err)
			}
		}
	}
}

// Fail : Override fail to require message
func (util *TestUtil) Fail(err ...interface{}) {
	util.Tester.Log(err...)
	util.Tester.Fail()
}

// FailE : Override fail to require messages "expected" and "got"
func (util *TestUtil) FailE(expected, got interface{}) {
	util.Tester.Log("Expected:", expected)
	util.Tester.Log("Got:", got)
	util.Tester.Fail()
}

// FailNow : Override fail now to require message
func (util *TestUtil) FailNow(err ...interface{}) {
	util.Tester.Log(err...)
	util.Tester.FailNow()
}

// UserInput : Apply user input to stdin
func (util *TestUtil) UserInput(input string) func() {
	tmpFile, err := ioutil.TempFile("", "NewUserInput")
	if err != nil {
		util.Fatal(err)
	}
	_, err = tmpFile.WriteString(input)
	if err != nil {
		util.Fatal(err)
	}
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		util.Fatal(err)
	}

	oldStdin := os.Stdin
	os.Stdin = tmpFile

	return func() {
		os.Stdin = oldStdin
		os.Remove(tmpFile.Name())
	}
}
