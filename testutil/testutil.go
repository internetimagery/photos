package testutil

import (
	"io/ioutil"
	"os"
	"testing"
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

// NewFile : Create a new file
func (util *TestUtil) NewFile(filePath string) {
	if err := ioutil.WriteFile(filePath, []byte("some info"), 0644); err != nil {
		util.Fatal(err)
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
func (util *TestUtil) Fail(err interface{}) {
	util.T.Log(err)
	util.T.Fail()
}

// FailE : Override fail to require messages "expected" and "got"
func (util *TestUtil) FailE(expected, got interface{}) {
	util.T.Log("Expected:", expected)
	util.T.Log("Got:", got)
	util.T.Fail()
}

// TempDir : Container for temporary directory
type TempDir struct {
	Dir string
	T   *testing.T
}

// NewTempDir : Create a new temporary directory
func NewTempDir(t *testing.T, prefix string) TempDir {
	tmpDir, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatal(err)
	}
	return TempDir{tmpDir, t}
}

// Close : Cleanup
func (tmp TempDir) Close() {
	err := os.RemoveAll(tmp.Dir)
	if err != nil {
		tmp.T.Fatal(err)
	}
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
