package copy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

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
	if err = ioutil.WriteFile(sourceFile, []byte("Testing 123"), 0644); err != nil {
		tu.Fatal(err)
	}

	// Test simple file copy works
	if err := <-File(sourceFile, destFile); err != nil {
		tu.Fail(err)
	}

	tu.AssertExists(destFile)
}
