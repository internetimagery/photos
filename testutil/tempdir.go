// Temporary Dir for testing.

package testutil

import (
	"io"
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

func (self *TempDir) Join(name string) string {
	return filepath.Join(self.Path, name)
}

func (self *TempDir) Copy(src string) string {
	in, err := os.Open(src)
	if err != nil {
		self.t.Fatal(err)
	}
	defer in.Close()

	dst := filepath.Join(self.Path, filepath.Base(src))
	out, err := os.Create(dst)
	if err != nil {
		self.t.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		self.t.Fatal(err)
	}
	defer out.Close()

	return dst
}

func (self *TempDir) Close() {
	err := os.RemoveAll(self.Path)
	if err != nil {
		self.t.Fatal(err)
	}
}
