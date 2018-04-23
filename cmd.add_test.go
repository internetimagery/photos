// Testing Add command

package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/internetimagery/photos/sandbox"
)

func TestTest(t *testing.T) {
	dir := sandbox.NewSandbox(t)
	defer dir.Close()
	// Do stuff
	files, _ := ioutil.ReadDir(dir.Path)
	fmt.Println(files[0].Name())
}
