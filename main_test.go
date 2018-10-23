package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/format"
	"github.com/internetimagery/photos/rename"
	"github.com/internetimagery/photos/testutil"
)

func TestQuestion(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.UserInput("y\n")()
	if !question() {
		tu.Fail("Question did not pass with 'y'")
	}

	defer tu.UserInput("n\n")
	if question() {
		tu.Fail("Question passed with 'n'")
	}

	defer tu.UserInput("anything\n")
	if question() {
		tu.Fail("Question passed with something other than 'y' / 'n'")
	}

}

// Test init
func TestInitClean(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Run init without a name
	defer tu.UserInput("y\n")()
	if err := run(tu.Dir, []string{"exe", "init"}); err == nil {
		tu.Fail("Allowed project with no name.")
	}

	// Run init on empty directory
	defer tu.UserInput("y\n")()
	if err := run(tu.Dir, []string{"exe", "init", "projectname"}); err != nil {
		tu.Fail(err)
	}

	// Ensure config file is created
	tu.AssertExists(filepath.Join(tu.Dir, context.ROOTCONF))
}

func TestInitExisting(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Run init on already set up directory
	defer tu.UserInput("y\n")()
	if err := run(tu.Dir, []string{"exe", "init", "projectname2"}); err == nil {
		tu.Fail("No error on already set up project.")
	}

	// Run in subfolder in setup directory
	subDir := filepath.Join(tu.Dir, "subdir")

	defer tu.UserInput("y\n")()
	if err := run(subDir, []string{"exe", "init", "projectname3"}); err == nil {
		tu.Fail("No error on already set up project in subfolder.")
	}
}

func TestSortClean(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")

	// Run sort on project not set up
	defer tu.UserInput("y\n")()
	if err := run(event, []string{"exe", "sort"}); err == nil {
		tu.Fail("Allowed usage on non-project folder.")
	}
}

func TestSortRoot(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Run sort on project not set up
	defer tu.UserInput("y\n")()
	if err := run(tu.Dir, []string{"exe", "sort"}); err == nil {
		tu.Fail("Allowed usage on root folder.")
	}
}

func TestSort(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	unsorted := filepath.Join(tu.Dir, "unsorted")

	// Run sort on root subdirectory
	defer tu.UserInput("y\n")()
	if err := run(unsorted, []string{"exe", "sort"}); err != nil {
		tu.Fail(err)
	}

	// Check files are where we expect them.
	tu.AssertExists(
		filepath.Join(unsorted, "18-10-23", "file1.txt"),
		filepath.Join(unsorted, "18-10-23", "file2.txt"),
		filepath.Join(unsorted, "18-10-23", "file2_1.txt"),
	)

}

// Testing rename command
func TestRename(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.TempDir("TestRename")()

	// Create an event
	subDir := filepath.Join(tu.Dir, "event01")
	tu.NewDir(subDir)

	// Make some test files
	testFiles := map[string]string{
		filepath.Join(subDir, "event01_002.test"):          filepath.Join(subDir, "event01_002.test"),
		filepath.Join(subDir, "event01_002[one two].test"): filepath.Join(subDir, "event01_002[one two].test"),
		filepath.Join(subDir, "newfile.test"):              filepath.Join(subDir, "event01_003.test"),
	}
	sourceTestFiles := []string{
		filepath.Join(subDir, rename.SOURCEDIR, "newfile.test"),
	}
	for testFile := range testFiles {
		tu.NewFile(testFile, "")
	}

	// Run without setting up project
	defer tu.UserInput("y\n")()
	if err := run(subDir, []string{"exe", "rename"}); err == nil {
		tu.Fail("Allowed running without project setup")
	}

	// Set up project
	tu.NewFile(filepath.Join(tu.Dir, context.ROOTCONF), `{
		"compress": [
			["*.missing", "cp -v \"$SOURCEPATH\" \"$DESTPATH.here\""]
		]}`)

	// Test rename in root folder
	defer tu.UserInput("y\n")()
	if err := run(tu.Dir, []string{"exe", "rename"}); err == nil {
		tu.Fail("Allowed running in root of project")
	}

	// Test rename
	defer tu.UserInput("y\n")()
	if err := run(subDir, []string{"exe", "rename"}); err != nil {
		tu.Fail(err)
	}

	// Check files are where they should be
	for _, testFile := range sourceTestFiles {
		tu.AssertExists(testFile)
	}
	for _, testFile := range testFiles {
		tu.AssertExists(testFile)
	}

	sourceDir := filepath.Join(subDir, rename.SOURCEDIR)
	testFiles = map[string]string{
		filepath.Join(sourceDir, "anotherfile.test"): filepath.Join(sourceDir, "anotherfile.test"),
		filepath.Join(subDir, "anotherfile.test"):    filepath.Join(sourceDir, "anotherfile_1.test"),
	}
	for testFile := range testFiles {
		tu.NewFile(testFile, "")
	}

	// Test rename again
	defer tu.UserInput("y\n")()
	if err := run(subDir, []string{"exe", "rename"}); err != nil {
		tu.Fail(err)
	}

	for _, testFile := range testFiles {
		tu.AssertExists(testFile)
	}

	// Test compress command is run
	tu.NewFile(filepath.Join(subDir, "testfile.missing"), "")
	expectFile := filepath.Join(subDir, format.TEMPPREFIX+"testfile.missing")

	// Expect rename to fail not finding compressed file
	defer tu.UserInput("y\n")()
	if err := run(subDir, []string{"exe", "rename"}); !os.IsNotExist(err) {
		if err == nil {
			tu.Fail("Did not alert failure to find compressed file.")
		} else {
			tu.Fail(err)
		}
	}

	// Check file was copied
	tu.AssertExists(expectFile)
}
