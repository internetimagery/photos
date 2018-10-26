package copy

import (
	"io/ioutil"
	"os"
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
	return &TestEnv{Dir: tmpDir, t: t}
}

func (env *TestEnv) Close() {
	if err := os.RemoveAll(env.Dir); err != nil {
		env.t.Fatal(err)
	}
}

func (env *TestEnv) MkDir(path string) string {
	if err := os.MkdirAll(path, 0755); err != nil {
		env.t.Fatal(err)
	}
	return path
}

func (env *TestEnv) MkFile(path string, data string, perm os.FileMode, modtime *time.Time) (string, os.FileInfo) {
	env.MkDir(filepath.Dir(path))
	if err := ioutil.WriteFile(path, []byte(data), perm); err != nil {
		env.t.Fatal(err)
	}
	if modtime != nil {
		if err := os.Chtimes(path, *modtime, *modtime); err != nil {
			env.t.Fatal(err)
		}
	}
	info, err := os.Stat(path)
	if err != nil {
		env.t.Fatal(err)
	}
	return path, info
}

func (env *TestEnv) Join(parts ...string) string {
	return filepath.Join(append([]string{env.Dir}, parts...)...)
}

func TestCopyFile(t *testing.T) {

	// Set up test environment (cannot use testutil.LoadTestdata() here)
	tu := NewTestEnv(t)
	defer tu.Close()

	perms := os.FileMode(0640)
	modtime := time.Date(2018, 10, 10, 0, 0, 0, 0, time.Local)
	sourceFile, sourceInfo := tu.MkFile(tu.Join("testfile1.txt"), "Hello", perms, &modtime)
	destFile := filepath.Join(tu.Dir, "testfile2.txt")

	// Record initial info
	initSize := sourceInfo.Size()

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
	if sourceInfoCompare.Size() != initSize || sourceInfoCompare.ModTime() != modtime {
		t.Log("Source file was modified!")
		t.Fail()
	}

	// Check our destination file matches the source
	destInfo, err := os.Stat(destFile)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if destInfo.ModTime() != modtime || destInfo.Size() != initSize {
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
	tu := NewTestEnv(t)
	defer tu.Close()

	sourceFile, sourceInfo := tu.MkFile(tu.Join("testfile1.txt"), "Hi", 0644, nil)
	destFile, destInfo := tu.MkFile(tu.Join("testfile2.txt"), "Hello", 0644, nil)

	sourceSize, sourceMod := sourceInfo.Size(), sourceInfo.ModTime()
	destSize, destMod := destInfo.Size(), destInfo.ModTime()

	// Test simple file copy works
	if err := <-File(sourceFile, destFile); !os.IsExist(err) {
		if err == nil {
			t.Log("No error with exsting file in place")
			t.Fail()
		} else {
			t.Log(err)
			t.Fail()
		}
	}

	// Check nothing changed
	sourceInfoCheck, err := os.Stat(sourceFile)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if sourceInfoCheck.Size() != sourceSize || sourceInfoCheck.ModTime() != sourceMod {
		t.Log("Source file has changed!")
		t.Fail()
	}
	destInfoCheck, err := os.Stat(destFile)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if destInfoCheck.Size() != destSize || destInfoCheck.ModTime() != destMod {
		t.Log("Destination file has changed!")
		t.Fail()
	}
}

func TestCopyFileNotFile(t *testing.T) {

	tu := NewTestEnv(t)
	defer tu.Close()

	sourceDir := tu.MkDir(tu.Join("testdir"))
	destDir := tu.Join("test2dir")

	// Test copy fails on non-files
	if err := <-File(sourceDir, destDir); err == nil {
		t.Log("No error with source directory")
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
