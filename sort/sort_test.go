package sort

import (
	"io/ioutil"
	"os"
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
	defer tu.TempDir("TestGetMediaDate")()

	testFile1 := filepath.Join(tu.Dir, "testfile.test")
	testTime1 := time.Now()
	if err := ioutil.WriteFile(testFile1, []byte("some stuff here"), 0644); err != nil {
		tu.Fatal(err)
	}

	compareTime, err := GetMediaDate(testFile1)
	if err != nil {
		tu.Fail(err)
	}

	layout := "06-01-02-15-04-05"

	if testTime1.Format(layout) != compareTime.Format(layout) {
		tu.FailE(testTime1, compareTime)
	}
}

func TestUniqueName(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.TempDir("TestUniqueName")()

	testFile1 := filepath.Join(tu.Dir, "test1.file") // File exists
	testFile2 := filepath.Join(tu.Dir, "test2.file") // File does not exist
	testExt := ".file"
	if err := ioutil.WriteFile(testFile1, []byte("stuff"), 0644); err != nil {
		tu.Fatal(err)
	}

	expectFile1 := filepath.Join(tu.Dir, "test1_1.file")
	compareFile1 := UniqueName(testFile1)
	compareFile2 := UniqueName(testFile2)
	compareExt := filepath.Ext(compareFile2)

	if compareFile1 != expectFile1 {
		tu.FailE(expectFile1, compareFile1)
	}
	if compareFile2 != testFile2 {
		tu.FailE(testFile2, compareFile2)
	}
	if compareExt != testExt {
		tu.FailE(testExt, compareExt)
	}

}

func TestSortMedia(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.TempDir("TestSortMedia")

	cxt := &context.Context{WorkingDir: tu.Dir}

	location, err := time.LoadLocation("")
	if err != nil {
		tu.Fatal(err)
	}
	modtime := time.Date(2018, 10, 16, 0, 0, 0, 0, location)
	folder := "18-10-16"

	testFiles := []string{
		filepath.Join(tu.Dir, "file1.txt"),
		filepath.Join(tu.Dir, "file2.txt"),
		filepath.Join(tu.Dir, folder, "file2.txt"),
	}

	if err = os.Mkdir(filepath.Join(tu.Dir, folder), 0755); err != nil {
		tu.Fatal(err)
	}
	for _, filename := range testFiles {
		if err = ioutil.WriteFile(filename, []byte("info"), 0644); err != nil {
			tu.Fatal(err)
		}
		if err = os.Chtimes(filename, modtime, modtime); err != nil {
			tu.Fatal(err)
		}
	}

	expectFiles := []string{
		filepath.Join(tmpDir.Dir, folder, "file1.txt"),
		filepath.Join(tmpDir.Dir, folder, "file2_1.txt"),
		filepath.Join(tmpDir.Dir, folder, "file2.txt"),
	}

	// Run our sort
	err = SortMedia(cxt)
	if err != nil {
		tu.Fail(err)
	}

	// Check our files made it to where they should be
	for _, file := range expectFiles {
		_, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				tu.Fail("Missing file", file)
			} else {
				tu.Fail(err)
			}

		}
	}

}
