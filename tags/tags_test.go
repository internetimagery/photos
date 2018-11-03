package tags

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/testutil"
)

// TODO: Add tag to untagged file
// TODO: Add tag, but duplicate tags in the function call
// TODO: Add tag to already tagged file
// TODO: Add tag to unformatted file
// TODO: Add duplicate tag

func TestAddTag(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Test adding a tag adjusts file
	testfile := filepath.Join(tu.Dir, "event01", "event01_001.txt")
	if err := AddTag(testfile, "one"); err != nil {
		tu.Fail(err)
	}
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001[one].txt"))
	// Test adding tag again
	if err := AddTag(testfile, "two"); err != nil {
		tu.Fail(err)
	}
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001[one two].txt"))
	// Test adding duplicate tag doesn't add it
	if err := AddTag(testfile, "two"); err != nil {
		tu.Fail(err)
	}
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001[one two].txt"))
	// Test adding duplicate tags and real tags still ignores duplicates
	if err := AddTag(testfile, "one", "two", "three"); err != nil {
		tu.Fail(err)
	}
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001[one two three].txt"))
	// Test adding no tags does nothing
	if err := AddTag(testfile, ""); err != nil {
		tu.Fail(err)
	}
	tu.AssertExists(filepath.Join(tu.Dir, "event01", "event01_001[one two three].txt"))
	// Test adding no tags
	testfile = filepath.Join(tu.Dir, "event01", "event01_002.txt")
	if err := AddTag(testfile, ""); err != nil {
		tu.Fail(err)
	}
	tu.AssertExists(testfile)
	// Test adding tags to unadded file does nothing
	testfile = filepath.Join(tu.Dir, "event01", "notpartofevent.txt")
	if err := AddTag(testfile, "one", "two"); err != nil {
		tu.Fail(err)
	}
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

}
