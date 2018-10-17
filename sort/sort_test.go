package sort

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	testDate := "08-10-16"
	location, err := time.LoadLocation("")
	if err != nil {
		t.Fatal(err)
	}
	testTime := time.Date(2008, 10, 16, 12, 0, 0, 0, location)

	compareDate := FormatDate(testTime)
	if compareDate != testDate {
		fmt.Println("Expected", testDate)
		fmt.Println("Got", compareDate)
		t.Fail()
	}
}

func TestGetMediaDate(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestGetMediaDate")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile1 := filepath.Join(tmpDir, "testfile.test")
	testTime1 := time.Now()
	if err = ioutil.WriteFile(testFile1, []byte("some stuff here"), 0644); err != nil {
		t.Fatal(err)
	}

	compareTime, err := GetMediaDate(testFile1)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if !(testTime1.Before(compareTime) && compareTime.After(testTime1)) {
		fmt.Println("Expected", testTime1)
		fmt.Println("Got", compareTime)
		t.Fail()
	}
}

func TestUniqueName(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestUniqueName")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile1 := filepath.Join(tmpDir, "test1.file")
	testFile2 := filepath.Join(tmpDir, "test2.file")
	testExt := ".file"
	if err = ioutil.WriteFile(testFile1, []byte("stuff"), 0644); err != nil {
		t.Fatal(err)
	}

	compareFile1 := UniqueName(testFile1)
	compareFile2 := UniqueName(testFile2)
	compareExt := filepath.Ext(compareFile2)

	if compareFile1 == testFile1 {
		fmt.Println("Filename not unique", testFile1)
		t.Fail()
	}
	if compareFile2 != testFile2 {
		fmt.Println("Expected", testFile2)
		fmt.Println("Got", compareFile2)
		t.Fail()
	}
	if compareExt != testExt {
		fmt.Println("Expected", testExt)
		fmt.Println("Got", compareExt)
		t.Fail()
	}

}
