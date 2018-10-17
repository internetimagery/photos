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
