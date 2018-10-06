// Sandbox utility for testing. Copy assets over to temporary test folder

package sandbox

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

// Types of media used in testing
const (
	CONFIGTYPE = iota
	IMAGETYPE
)

// AssetList : Hardcoded list of assets to copy over for testing. Mapped to types.
var AssetList = map[string]int{
	"photos.config":     CONFIGTYPE,
	"event01/img01.JPG": IMAGETYPE,
}

// assetRoot : Location of test assets from source path
var assetRoot = "sandbox/assets"

// SandBox : Temporary location, where tests can make or break things however they see fit
type SandBox struct {
	Root string     // Base temporary file housing the media
	t    *testing.T // Reference to the test, so we can fail if things break without having to return errors
}

// NewSandBox : Generate a new clean area to mess around with. Remember to defer "Close" to clean up.
func NewSandBox(t *testing.T) *SandBox {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(wd)
	_, root, _, _ := runtime.Caller(0)
	fmt.Println(root)
	return new(SandBox)
}

//
//
// // Copy assets to temp location for testing
// type SandBox struct {
// 	Path string
// 	t    *testing.T
// }

// func NewSandbox(t *testing.T) *SandBox {
// 	// Get source location. Testing always places working directory at project home.
// 	_, root, _, _ := runtime.Caller(0)
// 	root = filepath.Join(filepath.Dir(root), "assets")
//
// 	// Create temp dir
// 	tmp, err := ioutil.TempDir("", "Photos")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	// Copy files over to temp dir
// 	threads := make([]chan error, 0)
// 	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
// 		if err == nil && !info.IsDir() {
// 			base, err := filepath.Rel(root, path)
// 			if err != nil {
// 				return err
// 			}
// 			src, dst := filepath.Join(root, base), filepath.Join(tmp, base)
// 			done := make(chan error)
// 			threads = append(threads, done)
// 			go copy(src, dst, done)
// 		}
// 		return err
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	for _, done := range threads {
// 		err = <-done
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}
//
// 	return &SandBox{Path: tmp, t: t}
// }
//
// // Get asset at location
// func (self *SandBox) Get(name string) string {
// 	return filepath.Join(self.Path, name)
// }
//
// // Clean up
// func (self *SandBox) Close() {
// 	err := os.RemoveAll(self.Path)
// 	if err != nil {
// 		self.t.Fatal(err)
// 	}
// }
//
// // Simple file copy utility
// func copy(src, dst string, done chan error) {
// 	var err error
// 	defer func() {
// 		done <- err
// 	}()
// 	in, err := os.Open(src)
// 	if err != nil {
// 		return
// 	}
// 	defer in.Close()
//
// 	out, err := os.Create(dst)
// 	if err != nil {
// 		return
// 	}
// 	defer out.Close()
//
// 	_, err = io.Copy(out, in)
// 	if err != nil {
// 		return
// 	}
// 	defer out.Close()
// }
