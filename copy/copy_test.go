package copy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/internetimagery/photos/testutil"
)

func TestCopyFile(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	// Set up test environment (cannot use testutil.LoadTestdata() here)
	tmpDir, err := ioutil.TempDir("", "TestCopyFile")
	if err != nil {
		tu.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	sourceFile := filepath.Join(tmpDir, "testfile1.txt")
	destFile := filepath.Join(tmpDir, "testfile2.txt")
	perms := os.FileMode(0640)
	modtime := time.Date(2018, 10, 10, 0, 0, 0, 0, time.Local)
	data := []byte(strings.Repeat("TESTING AND", 100))
	fileSize := int64(len(data))
	if err = ioutil.WriteFile(sourceFile, data, perms); err != nil {
		tu.Fatal(err)
	}
	if err = os.Chtimes(sourceFile, modtime, modtime); err != nil {
		tu.Fatal(err)
	}

	// Test simple file copy works
	if err := <-File(sourceFile, destFile); err != nil {
		tu.Fail(err)
	}

	// Check everything matches
	info := tu.AssertExists(destFile)
	if info[0].ModTime() != modtime {
		tu.FailE(modtime, info[0].ModTime())
	}
	if info[0].Size() != fileSize {
		tu.Fail("File sizes vary")
	}
	// if runtime.GOOS != "windows" {
	// 	if info[0].Mode().Perm() != perms {
	// 		tu.FailE(perms, info[0].Mode().Perm())
	// 	}
	// }
}

func TestTree(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	// Set up test environment (cannot use testutil.LoadTestdata() here)
	tmpDir, err := ioutil.TempDir("", "TestCopyFile")
	if err != nil {
		tu.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

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
		tu.Fatal(err)
	}
	for _, testFile := range testFiles {
		if err = ioutil.WriteFile(testFile, []byte("info"), 0644); err != nil {
			tu.Fatal(err)
		}
	}

	// Do the thing
	if err = Tree(root1, root2); err != nil {
		tu.Fail(err)
	}

	// Check things exist!
	subDir = filepath.Join(root2, "subdir")
	tu.AssertExists(
		filepath.Join(root2, "testfile1.txt"),
		filepath.Join(root2, "testfile2.txt"),
		filepath.Join(subDir, "testfile3.txt"),
		filepath.Join(subDir, "testfile4.txt"),
	)
}
