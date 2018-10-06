package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSandbox(t *testing.T) {
	// Get the assets folder
	_, root, _, _ := runtime.Caller(0)
	assets := filepath.Join(filepath.Dir(root), "assets")

	// Create a new sandbox
	sb := NewSandBox(t)
	defer sb.Close()

	// Walk through the original assets and compare with the new sandbox version
	err := filepath.Walk(assets, func(sourcePath string, sourceInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !sourceInfo.IsDir() {
			// Get corresponding files
			relPath, _ := filepath.Rel(assets, sourcePath)
			destPath := filepath.Join(sb.Root, relPath)
			destInfo, err := os.Stat(destPath)
			if err != nil {
				return err
			}
			// Compare files
			if sourceInfo.Size() != destInfo.Size() {
				fmt.Println("Sizes are different", relPath)
				t.Fail()
			}
			if sourceInfo.ModTime() != destInfo.ModTime() {
				fmt.Println("Modification times are different", relPath)
				t.Fail()
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
