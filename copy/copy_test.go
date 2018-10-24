package copy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
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
	if err = ioutil.WriteFile(sourceFile, []byte("Testing 123"), perms); err != nil {
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
	if runtime.GOOS != "windows" {
		if info[0].Mode().Perm() != perms {
			tu.FailE(perms, info[0].Mode().Perm())
		}
	}
}
