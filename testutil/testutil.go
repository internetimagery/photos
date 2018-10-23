package testutil

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestUtil : Wrapper for helper functions
type TestUtil struct {
	Dir string
	*testing.T
}

// NewTestUtil : Create new testutil
func NewTestUtil(t *testing.T) *TestUtil {
	return &TestUtil{T: t}
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
	cmd := exec.Command("cp", "-avT", "--no-preserve", "ownership", testdata, tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		util.Log(string(output))
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

// NewFile : Create a new file
func (util *TestUtil) NewFile(filePath, content string) {
	if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
		util.Fatal(err)
	}
}

// NewDir : Create a new directory
func (util *TestUtil) NewDir(filePath string) {
	if err := os.MkdirAll(filePath, 0755); err != nil {
		util.Fatal(err)
	}
}

// AssertExists : Check if file exists. Fail if not
func (util *TestUtil) AssertExists(filePaths ...string) {
	for _, filePath := range filePaths {
		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				util.Fail("File does not exist:", filePath)
			} else {
				util.Fail(err)
			}
		}
	}
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

// TempDir : Create a new temporary directory
func (util *TestUtil) TempDir(name string) func() {
	if util.Dir != "" {
		util.Fatal("Tempfile already created", util.Dir)
	}
	tmpDir, err := ioutil.TempDir("", name)
	if err != nil {
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

// Fail : Override fail to require message
func (util *TestUtil) Fail(err ...interface{}) {
	util.T.Log(err...)
	util.T.Fail()
}

// FailE : Override fail to require messages "expected" and "got"
func (util *TestUtil) FailE(expected, got interface{}) {
	util.T.Log("Expected:", expected)
	util.T.Log("Got:", got)
	util.T.Fail()
}

// FailNow : Override fail now to require message
func (util *TestUtil) FailNow(err ...interface{}) {
	util.T.Log(err...)
	util.T.FailNow()
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
