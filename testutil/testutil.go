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
