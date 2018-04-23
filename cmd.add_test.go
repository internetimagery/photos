// Testing Add command

package main

import (
	"fmt"
	"testing"

	"github.com/internetimagery/photos/testutil"
)

func TestTest(t *testing.T) {
	dir := testutil.NewTempDir(t)
	defer dir.Close()
	// Do stuff
	fmt.Println(dir.Join("thing"))
}
