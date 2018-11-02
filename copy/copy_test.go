package copy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
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

////////////////////////////// TESTS /////////////////////////////

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

func TestCopyFileMissing(t *testing.T) {

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

	tu := NewTestEnv(t)
	defer tu.Close()

	sourceDir := tu.Join("root1")
	destDir := tu.Join("root2")

	tu.MkFile(tu.Join("root1", "test1.txt"), "", 0644, nil)
	tu.MkFile(tu.Join("root1", "subdir", "test2.txt"), "", 0644, nil)

	// Do the thing
	if err := Tree(sourceDir, destDir); err != nil {
		t.Log(err)
		t.Fail()
	}

	// Check files made it!
	if _, err := os.Stat(tu.Join("root2", "test1.txt")); err != nil {
		t.Log(err)
		t.Fail()
	}
	if _, err := os.Stat(tu.Join("root2", "subdir", "test2.txt")); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestTreeExistingDir(t *testing.T) {

	tu := NewTestEnv(t)
	defer tu.Close()

	sourceDir := tu.Join("root1")
	destDir := tu.Join("root2")

	tu.MkFile(tu.Join("root1", "test1.txt"), "", 0644, nil)
	tu.MkDir(destDir)

	// Do the thing
	if err := Tree(sourceDir, destDir); err != nil {
		t.Log(err)
		t.Fail()
	}

	// Check file made it!
	if _, err := os.Stat(tu.Join("root2", "test1.txt")); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestTreeExistingFile(t *testing.T) {

	tu := NewTestEnv(t)
	defer tu.Close()

	sourceDir := tu.Join("root1")
	destDir := tu.Join("root2")

	tu.MkFile(tu.Join("root1", "test1.txt"), "", 0644, nil)
	tu.MkFile(tu.Join("root1", "test2.txt"), "", 0644, nil)
	_, testInfo := tu.MkFile(tu.Join("root2", "test2.txt"), "Different", 0644, nil)

	testSize, testMod := testInfo.Size(), testInfo.ModTime()

	// Do the thing
	if err := Tree(sourceDir, destDir); !os.IsExist(err) {
		if err == nil {
			t.Log("File already exists and no error thrown")
			t.Fail()
		} else {
			t.Log(err)
			t.Fail()
		}
	}

	// Check file unchanged
	if _, err := os.Stat(tu.Join("root2", "test1.txt")); !os.IsNotExist(err) {
		if err != nil {
			t.Log("File copied across and not cleaned up", "root2/test1.txt")
		} else {
			t.Log(err)
		}
		t.Fail()
	}
	info, err := os.Stat(tu.Join("root2", "test2.txt"))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if info.Size() != testSize || info.ModTime() != testMod {
		t.Log("File was changed!")
		t.Fail()
	}
}

func TestTreeNoSourcePath(t *testing.T) {
	tu := NewTestEnv(t)
	defer tu.Close()

	sourcePath := tu.Join("one", "two", "three")
	destPath := tu.Join("four")

	if err := Tree(sourcePath, destPath); !os.IsNotExist(err) {
		if err == nil {
			t.Log("Failed to error on missing path")
		} else {
			t.Log(err)
		}
		t.Fail()
	}
}

func TestTreeSourceNotDir(t *testing.T) {
	tu := NewTestEnv(t)
	defer tu.Close()

	sourcePath, _ := tu.MkFile(tu.Join("one"), "here I am", 0644, nil)
	destPath := tu.Join("four")

	if err := Tree(sourcePath, destPath); err == nil {
		t.Log("Failed to error source file")
		t.Fail()
	}
}

func TestTreeDestNotDir(t *testing.T) {
	tu := NewTestEnv(t)
	defer tu.Close()

	sourcePath := tu.Join("one")
	destPath, _ := tu.MkFile(tu.Join("two"), "here I am", 0644, nil)

	if err := Tree(sourcePath, destPath); err == nil {
		t.Log("Failed to error dest file")
		t.Fail()
	}
}
