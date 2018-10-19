package testutil

import (
	"io/ioutil"
	"os"
	"testing"
)

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
func UserInput(t *testing.T, input string) func() {
	tmpFile, err := ioutil.TempFile("", "NewUserInput")
	if err != nil {
		t.Fatal(err)
	}
	_, err = tmpFile.WriteString(input)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	oldStdin := os.Stdin
	os.Stdin = tmpFile

	return func() {
		os.Stdin = oldStdin
		os.Remove(tmpFile.Name())
	}
}
