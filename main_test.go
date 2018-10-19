package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func TestQuestion(t *testing.T) {
	defer testutil.UserInput(t, "y\n")()
	if !question() {
		t.Log("Question did not pass with 'y'")
		t.Fail()
	}

	defer testutil.UserInput(t, "n\n")
	if question() {
		t.Log("Question passed with 'n'")
		t.Fail()
	}

}

// Test init
func TestInit(t *testing.T) {
	tmpDir := testutil.NewTempDir(t, "TestInitClean")
	defer tmpDir.Close()

	// Run init without a name
	defer testutil.UserInput(t, "y\n")()
	if err := run(tmpDir.Dir, []string{"exe", "init"}); err == nil {
		t.Log("Allowed project with no name.")
		t.Fail()
	}

	// Run init on empty directory
	defer testutil.UserInput(t, "y\n")()
	if err := run(tmpDir.Dir, []string{"exe", "init", "projectname"}); err != nil {
		t.Log(err)
		t.Fail()
	}

	// Ensure config file is created
	if _, err := os.Stat(filepath.Join(tmpDir.Dir, context.ROOTCONF)); err != nil {
		t.Log(err)
		t.Fail()
	}

	// Run init on already set up directory
	defer testutil.UserInput(t, "y\n")()
	if err := run(tmpDir.Dir, []string{"exe", "init", "projectname2"}); err == nil {
		t.Log("No error on already set up project.")
		t.Fail()
	}

	// Run in subfolder in setup directory
	subDir := filepath.Join(tmpDir.Dir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	defer testutil.UserInput(t, "y\n")()
	if err := run(subDir, []string{"exe", "init", "projectname3"}); err == nil {
		t.Log("No error on already set up project in subfolder.")
		t.Fail()
	}
}

// Test sort functionality
func TestSort(t *testing.T) {
	tmpDir := testutil.NewTempDir(t, "TestInitClean")
	defer tmpDir.Close()

	// Create subfolder
	subDir := filepath.Join(tmpDir.Dir, "subDir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Set up project
	defer testutil.UserInput(t, "y\n")()
	if err := run(tmpDir.Dir, []string{"exe", "init", "projectname"}); err != nil {
		t.Fatal(err)
	}

	// Set up test files
	testFolder1 := filepath.Join(subDir, "18-10-16")
	testFolder2 := filepath.Join(subDir, "18-10-17")
	loc, err := time.LoadLocation("")
	if err != nil {
		t.Fatal(err)
	}
	testDate1 := time.Date(2018, 10, 16, 0, 0, 0, 0, loc)
	testDate2 := time.Date(2018, 10, 17, 0, 0, 0, 0, loc)
	if err := os.Mkdir(testFolder2, 0755); err != nil {
		t.Fatal(err)
	}
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
		if err := ioutil.WriteFile(testFile.Test, []byte("info"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(testFile.Test, testFile.Date, testFile.Date); err != nil {
			t.Fatal(err)
		}
	}

	// Run sort on root directory
	defer testutil.UserInput(t, "y\n")()
	if err := run(tmpDir.Dir, []string{"exe", "sort"}); err == nil {
		t.Log("Allowing running sort in root... don't do that!")
		t.Fail()
	}

	// Run sort on root subdirectory
	defer testutil.UserInput(t, "y\n")()
	if err := run(subDir, []string{"exe", "sort"}); err != nil {
		t.Log(err)
		t.Fail()
	}

	// Check our files match!
	for _, testFile := range testFiles {
		if _, err := os.Stat(testFile.Expect); err != nil {
			t.Log(err)
			t.Fail()
		}
	}
}
