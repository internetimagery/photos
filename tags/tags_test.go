package tags

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/testutil"
)

func TestAddTag(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Test adding a tag adjusts file
	testfile := filepath.Join(tu.Dir, "event01", "18-11-26 event01_001.txt")
	tu.Must(AddTag([]string{testfile}, []string{"one"}))
	testfile = filepath.Join(tu.Dir, "event01", "18-11-26 event01_001[one].txt")
	tu.AssertExists(testfile)
	// Test adding tag again
	tu.Must(AddTag([]string{testfile}, []string{"two"}))
	testfile = filepath.Join(tu.Dir, "event01", "18-11-26 event01_001[one two].txt")
	tu.AssertExists(testfile)
	// Test adding duplicate tag doesn't add it
	tu.Must(AddTag([]string{testfile}, []string{"two"}))
	testfile = filepath.Join(tu.Dir, "event01", "18-11-26 event01_001[one two].txt")
	tu.AssertExists(testfile)
	// Test adding duplicate tags and real tags still ignores duplicates
	tu.Must(AddTag([]string{testfile}, []string{"one", "two", "three"}))
	testfile = filepath.Join(tu.Dir, "event01", "18-11-26 event01_001[one three two].txt")
	tu.AssertExists(testfile)
	// Test adding no tags does nothing
	tu.Must(AddTag([]string{testfile}, []string{""}))
	testfile = filepath.Join(tu.Dir, "event01", "18-11-26 event01_001[one three two].txt")
	tu.AssertExists(testfile)
	// Test adding no tags
	testfile = filepath.Join(tu.Dir, "event01", "18-11-26 event01_002.txt")
	tu.Must(AddTag([]string{testfile}, []string{""}))
	tu.AssertExists(testfile)
	// Test adding tags to unadded file does nothing
	testfile = filepath.Join(tu.Dir, "event01", "notpartofevent.txt")
	tu.Must(AddTag([]string{testfile}, []string{"one", "two"}))
	tu.AssertExists(testfile)
}

func TestAddTagExisting(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Test adding tag with existing file fails dramatically!
	testfile := filepath.Join(tu.Dir, "event01", "18-11-26 event01_001[one].txt")
	if err := AddTag([]string{testfile}, []string{"two"}); !os.IsExist(err) {
		if err == nil {
			tu.Fail("Succeeded in overwriting a file!")
		} else {
			tu.Fail(err)
		}
	}
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "18-11-26 event01_001[one].txt"))
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "18-11-26 event01_001[one two].txt"))
}

func TestRemoveTag(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Test removing tag from file with no tags does nothing
	testfile := filepath.Join(tu.Dir, "event01", "18-11-26 event01_001.txt")
	tu.Must(RemoveTag([]string{testfile}, []string{"one"}))
	tu.AssertExists(testfile)
	// Test removing tag from file removes tag... from file (and braces)
	testfile = filepath.Join(tu.Dir, "event01", "18-11-26 event01_002[one].txt")
	tu.Must(RemoveTag([]string{testfile}, []string{"one"}))
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "18-11-26 event01_002.txt"))
	// Test removing tag from file removes tag...
	testfile = filepath.Join(tu.Dir, "event01", "18-11-26 event01_003[one two].txt")
	tu.Must(RemoveTag([]string{testfile}, []string{"one"}))
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "18-11-26 event01_003[two].txt"))
	// Test removing tags that don't exist, does nothing
	testfile = filepath.Join(tu.Dir, "event01", "18-11-26 event01_004[one two].txt")
	tu.Must(RemoveTag([]string{testfile}, []string{"three"}))
	tu.AssertExists(testfile)
	// Test removing nothing does nothing
	testfile = filepath.Join(tu.Dir, "event01", "18-11-26 event01_004[one two].txt")
	tu.Must(RemoveTag([]string{testfile}, []string{""}))
	tu.AssertExists(testfile)
}

func TestRemoveTagExisting(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Test removing tag from file with no tags does nothing
	testfile := filepath.Join(tu.Dir, "event01", "18-11-26 event01_001[one two].txt")
	if err := RemoveTag([]string{testfile}, []string{"one"}); !os.IsExist(err) {
		if err == nil {
			tu.Fail("Allowed overwriting existing file!")
		} else {
			tu.Fail(err)
		}
	}
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "18-11-26 event01_001[two].txt"))
	tu.AssertExists(testfile)
}
