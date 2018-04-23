// Sandbox utility for testing. Copy assets over to temporary test folder

package sandbox

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Copy assets to temp location for testing
type SandBox struct {
	Path string
	t    *testing.T
}

func NewSandbox(t *testing.T) *SandBox {
	// Get source location
	_, root, _, _ := runtime.Caller(0)
	root = filepath.Join(filepath.Dir(root), "assets")

	// Create temp dir
	tmp, err := ioutil.TempDir("", "Photos")
	if err != nil {
		t.Fatal(err)
	}

	// Copy files over to temp dir
	threads := make([]chan error, 0)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			base, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			src, dst := filepath.Join(root, base), filepath.Join(tmp, base)
			done := make(chan error)
			threads = append(threads, done)
			go copy(src, dst, done)
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, done := range threads {
		err = <-done
		if err != nil {
			t.Fatal(err)
		}
	}

	return &SandBox{Path: tmp, t: t}
}

// Get asset at location
func (self *SandBox) Get(name string) string {
	return filepath.Join(self.Path, name)
}

// Clean up
func (self *SandBox) Close() {
	err := os.RemoveAll(self.Path)
	if err != nil {
		self.t.Fatal(err)
	}
}

// Simple file copy utility
func copy(src, dst string, done chan error) {
	var err error
	defer func() {
		done <- err
	}()
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}
	defer out.Close()
}
