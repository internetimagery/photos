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
	testfile := filepath.Join(tu.Dir, "event01", "event01_001.txt")
	tu.Must(AddTag(testfile, "one"))
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001[one].txt"))
	// Test adding tag again
	tu.Must(AddTag(testfile, "two"))
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001[one two].txt"))
	// Test adding duplicate tag doesn't add it
	tu.Must(AddTag(testfile, "two"))
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001[one two].txt"))
	// Test adding duplicate tags and real tags still ignores duplicates
	tu.Must(AddTag(testfile, "one", "two", "three"))
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001[one two three].txt"))
	// Test adding no tags does nothing
	tu.Must(AddTag(testfile, ""))
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001[one two three].txt"))
	// Test adding no tags
	testfile = filepath.Join(tu.Dir, "event01", "event01_002.txt")
	tu.Must(AddTag(testfile, ""))
	tu.AssertExists(testfile)
	// Test adding tags to unadded file does nothing
	testfile = filepath.Join(tu.Dir, "event01", "notpartofevent.txt")
	tu.Must(AddTag(testfile, "one", "two"))
	tu.AssertExists(testfile)
}

func TestAddTagExisting(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Test adding tag with existing file fails dramatically!
	testfile := filepath.Join(tu.Dir, "event01", "event01_001[one].txt")
	if err := AddTag(testfile, "two"); !os.IsExist(err) {
		if err == nil {
			tu.Fail("Succeeded in overwriting a file!")
		} else {
			tu.Fail(err)
		}
	}
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001[one].txt"))
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001[one two].txt"))
}

// TODO: Remove tag from tagged file
// TODO: Remove tag from file with multiple of the same tag
// TODO: Remove tag from file with multiple duplicates in the function call
// TODO: Remove tag from file that contains tags, but not the one in question
// TODO: Remove tag from file that contains no tags
// TODO: Remove tag from file that is not currently formatted

func TestRemoveTag(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Test removing tag from file with no tags does nothing
	testfile := filepath.Join(tu.Dir, "event01", "event01_001.txt")
	tu.Must(AddTag(testfile, "one"))
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001.txt"))
	// Test removing tag from file removes tag... from file (and braces)
	testfile = filepath.Join(tu.Dir, "event01", "event01_002[one].txt")
	tu.Must(AddTag(testfile, "one"))
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_002.txt"))
	// Test removing tag from file removes tag...
	testfile = filepath.Join(tu.Dir, "event01", "event01_003[one two].txt")
	tu.Must(AddTag(testfile, "one"))
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_003[two].txt"))
	// Test removing tags that don't exist, does nothing
	testfile = filepath.Join(tu.Dir, "event01", "event01_004[one two].txt")
	tu.Must(AddTag(testfile, "three"))
	tu.AssertExists(testfile)
	// Test removing nothing does nothing
	testfile = filepath.Join(tu.Dir, "event01", "event01_004[one two].txt")
	tu.Must(AddTag(testfile, ""))
	tu.AssertExists(testfile)
}
