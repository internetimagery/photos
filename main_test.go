package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func TestQuestion(t *testing.T) {
	defer testutil.UserInput(t, "y\n")()
	if !question() {
		fmt.Println("Question did not pass with 'y'")
		t.Fail()
	}

	defer testutil.UserInput(t, "n\n")
	if question() {
		fmt.Println("Question passed with 'n'")
		t.Fail()
	}

}

// Test init
func TestInit(t *testing.T) {
	tmpDir := testutil.NewTempDir(t, "TestInitClean")
	defer tmpDir.Close()

	// Run init on empty directory
	defer testutil.UserInput(t, "y\n")()
	if err := run(tmpDir.Dir, []string{"exe", "init", "projectname"}); err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Ensure config file is created
	if _, err := os.Stat(filepath.Join(tmpDir.Dir, context.ROOTCONF)); err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Run init on already set up directory
	defer testutil.UserInput(t, "y\n")()
	if err := run(tmpDir.Dir, []string{"exe", "init", "projectname2"}); err == nil {
		fmt.Println("No error on already set up project.")
		t.Fail()
	}

	// Run in subfolder in setup directory
	subDir := filepath.Join(tmpDir.Dir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	defer testutil.UserInput(t, "y\n")()
	if err := run(subDir, []string{"exe", "init", "projectname3"}); err == nil {
		fmt.Println("No error on already set up project in subfolder.")
		t.Fail()
	}
}
