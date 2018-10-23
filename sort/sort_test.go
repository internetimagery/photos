package sort

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func TestFormatDate(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	testDate := "08-10-16"
	location, err := time.LoadLocation("")
	if err != nil {
		tu.Fatal(err)
	}
	testTime := time.Date(2008, 10, 16, 12, 0, 0, 0, location)

	compareDate := FormatDate(testTime)
	if compareDate != testDate {
		tu.FailE(testDate, compareDate)
	}
}

func TestGetMediaDate(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	testFile1 := filepath.Join(tu.Dir, "testfile.txt")
	testTime1 := "18-10-22"
	tu.ModTime(2018, 10, 22, testFile1)

	compareTime, err := GetMediaDate(testFile1)
	if err != nil {
		tu.Fail(err)
	}

	layout := "06-01-02"

	if testTime1 != compareTime.Format(layout) {
		tu.FailE(testTime1, compareTime)
	}
}

func TestUniqueName(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	testFile1 := filepath.Join(tu.Dir, "test1.file") // File exists
	testFile2 := filepath.Join(tu.Dir, "test2.file") // File does not exist
	testExt := ".file"

	expectFile1 := filepath.Join(tu.Dir, "test1_1.file")
	expectFile2 := filepath.Join(tu.Dir, "test2.file")

	compareFile1 := UniqueName(testFile1)
	compareFile2 := UniqueName(testFile2)
	compareExt := filepath.Ext(compareFile2)

	if compareFile1 != expectFile1 {
		tu.FailE(expectFile1, compareFile1)
	}
	if compareFile2 != expectFile2 {
		tu.FailE(expectFile2, compareFile2)
	}
	if compareExt != testExt {
		tu.FailE(testExt, compareExt)
	}
}

func TestSortMedia(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	// Get our context
	cxt, err := context.NewContext(tu.Dir)
	if err != nil {
		tu.Fatal(err)
	}

	dateDir := filepath.Join(tu.Dir, "18-10-22")
	tu.ModTime(2018, 10, 22,
		filepath.Join(tu.Dir, "file1.txt"),
		filepath.Join(tu.Dir, "file2.txt"),
	)

	// Run our sort
	err = SortMedia(cxt)
	if err != nil {
		tu.Fail(err)
	}

	tu.AssertExists(
		filepath.Join(dateDir, "file1.txt"),
		filepath.Join(dateDir, "file2.txt"),
		filepath.Join(dateDir, "file2_1.txt"),
	)
}
