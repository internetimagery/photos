package sort

import (
	"fmt"
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
		tu.Fail(fmt.Sprintf("Expected '%s'\nGot '%s'", testDate, compareDate))
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
		tu.Fail(fmt.Sprintf("Expected '%s'\nGot '%s'", testTime1, compareTime))
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
		tu.Fail(fmt.Sprintf("Expected '%s'\nGot '%s'", expectFile1, compareFile1))
	}
	if compareFile2 != testFile2 {
		tu.Fail(fmt.Sprintf("Expected '%s'\nGot '%s'", testFile2, compareFile2))
	}
	if compareExt != testExt {
		fmt.Println("Expected", testExt)
		fmt.Println("Got", compareExt)
		t.Fail()
	}

}

func TestSortMedia(t *testing.T) {
	tmpDir := testutil.NewTempDir(t, "TestSortMedia")
	defer tmpDir.Close()

	cxt := &context.Context{WorkingDir: tmpDir.Dir}

	location, err := time.LoadLocation("")
	if err != nil {
		t.Fatal(err)
	}
	modtime := time.Date(2018, 10, 16, 0, 0, 0, 0, location)
	folder := "18-10-16"

	testFiles := []string{
		filepath.Join(tmpDir.Dir, "file1.txt"),
		filepath.Join(tmpDir.Dir, "file2.txt"),
		filepath.Join(tmpDir.Dir, folder, "file2.txt"),
	}

	if err = os.Mkdir(filepath.Join(tmpDir.Dir, folder), 0755); err != nil {
		t.Fatal(err)
	}
	for _, filename := range testFiles {
		if err = ioutil.WriteFile(filename, []byte("info"), 0644); err != nil {
			t.Fatal(err)
		}
		if err = os.Chtimes(filename, modtime, modtime); err != nil {
			t.Fatal(err)
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
		fmt.Println(err)
		t.Fail()
	}

	// Check our files made it to where they should be
	for _, file := range expectFiles {
		_, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Missing file", file)
				t.Fail()
			} else {
				fmt.Println(err)
				t.Fail()
			}

		}
	}

}
