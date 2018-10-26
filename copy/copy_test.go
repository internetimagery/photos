package copy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCopyFile(t *testing.T) {

	// Set up test environment (cannot use testutil.LoadTestdata() here)
	tmpDir, err := ioutil.TempDir("", "TestCopyFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	sourceFile := filepath.Join(tmpDir, "testfile1.txt")
	destFile := filepath.Join(tmpDir, "testfile2.txt")
	perms := os.FileMode(0640)
	modtime := time.Date(2018, 10, 10, 0, 0, 0, 0, time.Local)
	data := []byte(strings.Repeat("TESTING AND", 100))
	fileSize := int64(len(data))
	if err = ioutil.WriteFile(sourceFile, data, perms); err != nil {
		t.Fatal(err)
	}
	if err = os.Chtimes(sourceFile, modtime, modtime); err != nil {
		t.Fatal(err)
	}

	// Test simple file copy works
	if err := <-File(sourceFile, destFile); err != nil {
		t.Log(err)
		t.Fail()
	}

	// Check everything matches
	info, err := os.Stat(destFile)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if info.ModTime() != modtime {
		t.Log("Expected", modtime)
		t.Log("Got", info.ModTime())
		t.Fail()
	}
	if info.Size() != fileSize {
		t.Log("File sizes vary")
		t.Fail()
	}
	// if runtime.GOOS != "windows" {
	// 	if info[0].Mode().Perm() != perms {
	// 		tu.FailE(perms, info[0].Mode().Perm())
	// 	}
	// }
}

func TestTree(t *testing.T) {

	// Set up test environment (cannot use testutil.LoadTestdata() here)
	tmpDir, err := ioutil.TempDir("", "TestCopyFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testData := []byte(strings.Repeat("TESTING AGAIN", 100))
	testDataSize := int64(len(testData))
	root1 := filepath.Join(tmpDir, "root1")
	root2 := filepath.Join(tmpDir, "root2")

	// Setting up test env
	subDir := filepath.Join(root1, "subdir")
	testFiles := []string{
		filepath.Join(root1, "testfile1.txt"),
		filepath.Join(root1, "testfile2.txt"),
		filepath.Join(subDir, "testfile3.txt"),
		filepath.Join(subDir, "testfile4.txt"),
	}
	if err = os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	for _, testFile := range testFiles {
		if err = ioutil.WriteFile(testFile, testData, 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Do the thing
	if err = Tree(root1, root2); err != nil {
		t.Log(err)
		t.Fail()
	}

	// Check things exist!
	subDir = filepath.Join(root2, "subdir")
	resultFiles := []string{
		filepath.Join(root2, "testfile1.txt"),
		filepath.Join(root2, "testfile2.txt"),
		filepath.Join(subDir, "testfile3.txt"),
		filepath.Join(subDir, "testfile4.txt"),
	}
	for _, resultFile := range resultFiles {
		info, err := os.Stat(resultFile)
		if err != nil {
			t.Log(err)
			t.Fail()
		} else if info.Size() != testDataSize {
			t.Log("File sizes differ", resultFile)
			t.Fail()
		}
	}
}
