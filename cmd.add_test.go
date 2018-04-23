// Testing Add command

package main

import (
	"os"
	"testing"

	"github.com/internetimagery/photos/sandbox"
)

func TestTest(t *testing.T) {
	dir := sandbox.NewSandbox(t)
	defer dir.Close()
	// Check asset is there
	asset := dir.Get("img1.jpg")
	if _, err := os.Stat(asset); os.IsNotExist(err) {
		t.Fail()
	}
}
