package copy

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

type TestEnv struct {
	Dir string
	t   *testing.T
}

func NewTestEnv(t *testing.T) *TestEnv {
	tmpDir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("cp", "avT", filepath.Join("testdata", t.Name()), tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Log(string(output))
		t.Fatal(err)
	}
	return &TestEnv{Dir: tmpDir, t: t}
}

func (env *TestEnv) Close() {
	if err := os.RemoveAll(env.Dir); err != nil {
		env.t.Fatal(err)
	}
}

func TestCopyFile(t *testing.T) {

	// Set up test environment (cannot use testutil.LoadTestdata() here)
	tu := NewTestEnv(t)
	defer tu.Close()

	sourceFile := filepath.Join(tu.Dir, "testfile1.txt")
	destFile := filepath.Join(tu.Dir, "testfile2.txt")
	perms := os.FileMode(0640)
	modtime := time.Date(2018, 10, 10, 0, 0, 0, 0, time.Local)
	if err := os.Chtimes(sourceFile, modtime, modtime); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(sourceFile, perms); err != nil {
		t.Fatal(err)
	}

	// Grab info about file
	sourceInfo, err := os.Stat(sourceFile)
	if err != nil {
		t.Fatal(err)
	}

	// Test simple file copy works
	if err := <-File(sourceFile, destFile); err != nil {
		t.Log(err)
		t.Fail()
	}

	// Check our source file remains untouched
	sourceInfoCompare, err := os.Stat(sourceFile)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if sourceInfoCompare.Size() != sourceInfo.Size() || sourceInfoCompare.ModTime() != modtime {
		t.Log("Source file was modified!")
		t.Fail()
	}

	// Check our destination file matches the source
	destInfo, err := os.Stat(destFile)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if destInfo.ModTime() != modtime || destInfo.Size() != sourceInfo.Size() {
		t.Log("Expected", sourceInfo)
		t.Log("Got", destInfo)
		t.Fail()
	}
	if runtime.GOOS != "windows" {
		if destInfo.Mode().Perm() != perms {
			t.Log("Expected", perms)
			t.Log("Got", destInfo.Mode().Perm())
			t.Fail()
		}
	}
}

func TestCopyFileExisting(t *testing.T) {

	// Set up test environment (cannot use testutil.LoadTestdata() here)
	tmpDir, err := ioutil.TempDir("", "TestCopyFileExisting")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	sourceFile := filepath.Join(tmpDir, "testfile1.txt")
	destFile := filepath.Join(tmpDir, "testfile2.txt")
	sourceData := []byte(strings.Repeat("TESTING AND", 100))
	destData := []byte(strings.Repeat("RESULT AND", 100))
	if err = ioutil.WriteFile(sourceFile, sourceData, 0644); err != nil {
		t.Fatal(err)
	}
	if err = ioutil.WriteFile(destFile, destData, 0644); err != nil {
		t.Fatal(err)
	}

	// Test simple file copy works
	if err := <-File(sourceFile, destFile); !os.IsExist(err) {
		if err == nil {
			t.Log("No error with exsting file")
			t.Fail()
		} else {
			t.Log(err)
			t.Fail()
		}
	}

	// Check nothing changed
	data, err := ioutil.ReadFile(destFile)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if !bytes.Equal(data, destData) {
		t.Log("Data was changed!")
		t.Fail()
	}
}

func TestCopyFileNotFile(t *testing.T) {

	// Set up test environment (cannot use testutil.LoadTestdata() here)
	tmpDir, err := ioutil.TempDir("", "TestCopyFileNotFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	sourceFile := filepath.Join(tmpDir, "testfile1.txt")
	destFile := filepath.Join(tmpDir, "testfile2.txt")
	if err = ioutil.WriteFile(sourceFile, sourceData, 0644); err != nil {
		t.Fatal(err)
	}
	if err = ioutil.WriteFile(destFile, destData, 0644); err != nil {
		t.Fatal(err)
	}

	// Test simple file copy works
	if err := <-File(sourceFile, destFile); !os.IsExist(err) {
		if err == nil {
			t.Log("No error with exsting file")
			t.Fail()
		} else {
			t.Log(err)
			t.Fail()
		}
	}

	// Check nothing changed
	data, err := ioutil.ReadFile(destFile)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if !bytes.Equal(data, destData) {
		t.Log("Data was changed!")
		t.Fail()
	}
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
