package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/internetimagery/photos/context"
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

}

// Test init
func TestInit(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.TempDir("TestInit")()

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

	// Run init on already set up directory
	defer tu.UserInput("y\n")()
	if err := run(tu.Dir, []string{"exe", "init", "projectname2"}); err == nil {
		tu.Fail("No error on already set up project.")
	}

	// Run in subfolder in setup directory
	subDir := filepath.Join(tu.Dir, "subdir")
	tu.NewDir(subDir)

	defer tu.UserInput("y\n")()
	if err := run(subDir, []string{"exe", "init", "projectname3"}); err == nil {
		tu.Fail("No error on already set up project in subfolder.")
	}
}

// Test sort functionality
func TestSort(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.TempDir("TestSort")()

	// Create subfolder
	subDir := filepath.Join(tu.Dir, "subDir")
	tu.NewDir(subDir)

	// Run sort on project not set up
	defer tu.UserInput("y\n")()
	if err := run(tu.Dir, []string{"exe", "sort"}); err == nil {
		tu.Fail("Allowed usage on non-project folder.")
	}

	// Set up project
	defer tu.UserInput("y\n")()
	if err := run(tu.Dir, []string{"exe", "init", "projectname"}); err != nil {
		tu.Fatal(err)
	}

	// Set up test files
	testFolder1 := filepath.Join(subDir, "18-10-16")
	testFolder2 := filepath.Join(subDir, "18-10-17")
	loc, err := time.LoadLocation("")
	if err != nil {
		tu.Fatal(err)
	}
	testDate1 := time.Date(2018, 10, 16, 0, 0, 0, 0, loc)
	testDate2 := time.Date(2018, 10, 17, 0, 0, 0, 0, loc)
	tu.NewDir(testFolder2)

	type testCase struct {
		Test, Expect string
		Date         time.Time
	}
	testFiles := []testCase{
		testCase{Test: filepath.Join(subDir, "file1.txt"), Expect: filepath.Join(testFolder1, "file1.txt"), Date: testDate1},      // Standard file
		testCase{Test: filepath.Join(subDir, "file2.txt"), Expect: filepath.Join(testFolder2, "file2_1.txt"), Date: testDate2},    // Second file, moddate in testFolder
		testCase{Test: filepath.Join(testFolder2, "file2.txt"), Expect: filepath.Join(testFolder2, "file2.txt"), Date: testDate2}, // File of same name
	}
	for _, testFile := range testFiles {
		tu.NewFile(testFile.Test, "")
		if err := os.Chtimes(testFile.Test, testFile.Date, testFile.Date); err != nil {
			tu.Fatal(err)
		}
	}

	// Run sort on root directory
	defer tu.UserInput("y\n")()
	if err := run(tu.Dir, []string{"exe", "sort"}); err == nil {
		tu.Fail("Allowing running sort in root... don't do that!")
	}

	// Run sort on root subdirectory
	defer tu.UserInput("y\n")()
	if err := run(subDir, []string{"exe", "sort"}); err != nil {
		tu.Fail(err)
	}

	// Check our files match!
	for _, testFile := range testFiles {
		tu.AssertExists(testFile.Expect)
	}
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
	defer tu.UserInput("y\n")()
	if err := run(tu.Dir, []string{"exe", "init", "projectname"}); err != nil {
		tu.Fatal(err)
	}

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

}
