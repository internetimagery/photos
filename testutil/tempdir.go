// Temporary Dir for testing.

package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type TempDir struct {
	Path string
	t    *testing.T
}

func NewTempDir(t *testing.T) *TempDir {
	dir, err := ioutil.TempDir("", "Photos")
	if err != nil {
		t.Fatal(err)
	}
	return &TempDir{Path: dir, t: t}
}

func (self *TempDir) Add(name string) string {
	return filepath.Join(self.Path, name)
}

func (self *TempDir) Close() {
	err := os.RemoveAll(self.Path)
	if err != nil {
		self.t.Fatal(err)
	}
}
