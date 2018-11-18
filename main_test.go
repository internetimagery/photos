package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/lock"
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
	tu.Must(run(tu.Dir, []string{"exe", "init", "projectname"}))

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
	if err := run(event, []string{"exe", "sort", "event01"}); err == nil {
		tu.Fail("Allowed usage on non-project folder.")
	}
}

func TestSortRoot(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Run sort on project not set up
	defer tu.UserInput("y\n")()
	if err := run(tu.Dir, []string{"exe", "sort", "."}); err == nil {
		tu.Fail("Allowed usage on root folder.")
	}
}

func TestSort(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	project := filepath.Join(tu.Dir, "project")
	sorted := filepath.Join(project, "Sorted")

	tu.ModTime(2018, 10, 10, filepath.Join(tu.Dir, "file1.txt"))
	tu.ModTime(2018, 10, 23, filepath.Join(tu.Dir, "file2.txt"))

	// Run sort without any input
	defer tu.UserInput("y\n")()
	if err := run(project, []string{"exe", "sort"}); err == nil {
		tu.Fail("Allowed no source input")
	}

	// Run sort on files
	defer tu.UserInput("y\n")()
	tu.Must(run(project, []string{"exe", "sort", "../"}))

	// Check files are where we expect them.
	tu.AssertExists(
		filepath.Join(sorted, "18-10-10", "file1.txt"),
		filepath.Join(sorted, "18-10-23", "file2.txt"),
		filepath.Join(sorted, "18-10-23", "file2_1.txt"),
	)
}

func TestRenameClean(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")

	// Run without setting up project
	defer tu.UserInput("y\n")()
	if err := run(event, []string{"exe", "rename"}); err == nil {
		tu.Fail("Allowed running without project setup")
	}
}

func TestRenameRoot(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Run without setting up project
	defer tu.UserInput("y\n")()
	if err := run(tu.Dir, []string{"exe", "rename"}); err == nil {
		tu.Fail("Allowed running in root")
	}
}

func TestRenameMissing(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")

	// Test rename failing if compress command fails to produce output
	defer tu.UserInput("y\n")()
	if err := run(event, []string{"exe", "rename"}); !os.IsNotExist(err) {
		if err == nil {
			tu.Fail("Did not alert failure to find compressed file.")
		} else {
			tu.Fail(err)
		}
	}

	// Check command actually put our file there.
	tu.AssertExists(
		filepath.Join(event, "tmp-testfile.missing.here"),
	)
}

func TestRenameExistingSource(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")

	// Test rename with existing file in source dir
	defer tu.UserInput("y\n")()
	tu.Must(run(event, []string{"exe", "rename"}))

	tu.AssertExists(
		filepath.Join(event, rename.SOURCEDIR, "testfile1.txt"),
		filepath.Join(event, rename.SOURCEDIR, "testfile2.txt"),
		filepath.Join(event, rename.SOURCEDIR, "testfile2_1.txt"),
	)
}

// Testing rename command
func TestRename(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")

	// Test rename
	defer tu.UserInput("y\n")()
	tu.Must(run(event, []string{"exe", "rename"}))

	tu.AssertExists(
		filepath.Join(event, "event01_002.test"),
		filepath.Join(event, "event01_002[one two].test"),
		filepath.Join(event, "event01_003.test"),
		filepath.Join(event, rename.SOURCEDIR, "newfile.test"),
	)
}

func TestLock(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")
	if err, ok := run(event, []string{"exe", "lock"}).(*lock.MissmatchError); !ok {
		if err == nil {
			tu.Fail("Allowed missmatch lock")
		} else {
			tu.Fail(err)
		}
	}
	tu.Must(run(event, []string{"exe", "lock", "--force"}))
}

func TestLockRoot(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	if err := run(tu.Dir, []string{"exe", "lock"}); err == nil {
		tu.Fail("Allowed locking root")
	}
}

func TestAddTag(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Add tag by filename
	event := filepath.Join(tu.Dir, "event01")
	tu.Must(run(event, []string{"exe", "tag", filepath.Join(event, "event01_010[one].txt"), "two", "three"}))
	tu.AssertExists(filepath.Join(event, "event01_010[one three two].txt"))

	// Add tag by index
	tu.Must(run(event, []string{"exe", "tag", "10", "four"}))
	tu.AssertExists(filepath.Join(event, "event01_010[four one three two].txt"))

	// Add numeric (index looking) tag by index using --
	tu.Must(run(event, []string{"exe", "tag", "10", "--", "5"}))
	tu.AssertExists(filepath.Join(event, "event01_010[5 four one three two].txt"))

	// no tags, stopping with --
	if err := run(event, []string{"exe", "tag", "10", "--"}); err == nil {
		tu.Fail("Allowed no tags to be specified.")
	}

	// no files, starting with --
	if err := run(event, []string{"exe", "tag", "--", "10"}); err == nil {
		tu.Fail("Allowed no files to be specified.")
	}

}

func TestRemoveTag(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Remove tag by filename
	event := filepath.Join(tu.Dir, "event01")
	tu.Must(run(event, []string{"exe", "tag", "--remove", filepath.Join(event, "event01_010[5 three one two].txt"), "two"}))
	tu.AssertExists(filepath.Join(event, "event01_010[5 one three].txt"))

	// Remove tag by index
	tu.Must(run(event, []string{"exe", "tag", "--remove", "10", "one"}))
	tu.AssertExists(filepath.Join(event, "event01_010[5 three].txt"))

	// Remove numeric (index looking) tag by index using --
	tu.Must(run(event, []string{"exe", "tag", "--remove", "10", "--", "5"}))
	tu.AssertExists(filepath.Join(event, "event01_010[three].txt"))

	// --remove in different spot
	tu.Must(run(event, []string{"exe", "tag", "10", "--remove", "three"}))
	tu.AssertExists(filepath.Join(event, "event01_010.txt"))

	// --remove in tag section
	tu.Must(run(event, []string{"exe", "tag", "10", "--", "--remove", "three"}))
	tu.AssertExists(filepath.Join(event, "event01_010.txt"))

	// no tags, stopping with --
	if err := run(event, []string{"exe", "tag", "10", "--", "--remove"}); err == nil {
		tu.Fail("Allowed no tags to be specified.")
	}

	// no files, starting with --
	if err := run(event, []string{"exe", "tag", "--remove", "--", "10"}); err == nil {
		tu.Fail("Allowed no files to be specified.")
	}

}
