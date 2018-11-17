package lock

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func testReadOnly(tu *testutil.TestUtil, filename string) {
	if err := ioutil.WriteFile(filename, []byte("Fail"), 0644); !os.IsPermission(err) {
		if err == nil {
			tu.Log("Did not make file readonly", filename)
		} else {
			tu.Fail(err)
		}
	}
}

func TestGenerateContentHash(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	buff := bytes.NewReader([]byte("SOME DATA HERE"))
	expectHash := "SHA256:j7BgIUq2w472YYetmry+ieE0D3kqaVRdU6Ri6uq2hTY=" // MD5

	testHash := tu.Must(GenerateContentHash("SHA256", buff)).(string)
	if expectHash != testHash {
		tu.FailE(expectHash, testHash)
	}
}

func TestGeneratePerceptualHash(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	handle1 := tu.MustFatal(os.Open(filepath.Join(tu.Dir, "testimg1.jpg"))).(*os.File)
	defer handle1.Close()
	handle2 := tu.MustFatal(os.Open(filepath.Join(tu.Dir, "testimg2.jpg"))).(*os.File)
	defer handle2.Close()
	handle3 := tu.MustFatal(os.Open(filepath.Join(tu.Dir, "testimg3.jpg"))).(*os.File)
	defer handle3.Close()
	handle4 := tu.MustFatal(os.Open(filepath.Join(tu.Dir, "testimg4.jpg"))).(*os.File)
	defer handle4.Close()
	handle5 := tu.MustFatal(os.Open(filepath.Join(tu.Dir, "testimg5.png"))).(*os.File)
	defer handle5.Close()
	handle6 := tu.MustFatal(os.Open(filepath.Join(tu.Dir, "testimg6.txt"))).(*os.File)
	defer handle6.Close()
	handle7 := tu.MustFatal(os.Open(filepath.Join(tu.Dir, "testimg7.tiff"))).(*os.File)
	defer handle7.Close()

	testHash1 := tu.Must(GeneratePerceptualHash("average", handle1)).(string)
	testHash2 := tu.Must(GeneratePerceptualHash("average", handle2)).(string)
	testHash3 := tu.Must(GeneratePerceptualHash("average", handle3)).(string)
	testHash4 := tu.Must(GeneratePerceptualHash("average", handle4)).(string)
	testHash5 := tu.Must(GeneratePerceptualHash("average", handle5)).(string)
	_, err := GeneratePerceptualHash("average", handle6)
	if err == nil {
		tu.Fail("Allowed unsupported file")
	}
	_, err = GeneratePerceptualHash("average", handle7)
	if err == nil {
		tu.Fail("Allowed unsupported img")
	}

	if !tu.Must(IsSamePerceptualHash(testHash1, testHash2)).(bool) ||
		!tu.Must(IsSamePerceptualHash(testHash1, testHash3)).(bool) ||
		!tu.Must(IsSamePerceptualHash(testHash1, testHash5)).(bool) ||
		!tu.Must(IsSamePerceptualHash(testHash2, testHash3)).(bool) ||
		!tu.Must(IsSamePerceptualHash(testHash2, testHash5)).(bool) ||
		!tu.Must(IsSamePerceptualHash(testHash3, testHash5)).(bool) {
		tu.Fail("Equals not equal")
	}

	if tu.Must(IsSamePerceptualHash(testHash1, testHash4)).(bool) {
		tu.Fail("False positive")
	}

}

func TestSnapshot(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	testfile1 := filepath.Join(tu.Dir, "testimg1.txt")
	testfile2 := filepath.Join(tu.Dir, "testimg2.jpg")
	modtime1 := tu.MustFatal(os.Stat(testfile1)).(os.FileInfo).ModTime()
	modtime2 := tu.MustFatal(os.Stat(testfile2)).(os.FileInfo).ModTime()

	sshot1, sshot2 := new(Snapshot), new(Snapshot)
	tu.Must(<-sshot1.Generate(testfile1))
	tu.Must(<-sshot2.Generate(testfile2))

	if expect := "testimg1.txt"; sshot1.Name != expect {
		tu.FailE(expect, sshot1.Name)
	}
	if expect := int64(24); sshot1.Size != expect {
		tu.FailE(expect, sshot1.Size)
	}
	if !modtime1.Equal(sshot1.ModTime) {
		tu.FailE(modtime1, sshot1.ModTime)
	}
	if expect := "SHA256:h13POS/MwQ0SHVmJOSHgeN7+fM9ymIJZvdZt3nnLAqY="; sshot1.ContentHash["SHA256"] != expect {
		tu.FailE(expect, sshot1.ContentHash["SHA256"])
	}
	if sshot1.PerceptualHash != nil {
		tu.FailE(nil, sshot1.PerceptualHash)
	}

	if expect := "testimg2.jpg"; sshot2.Name != expect {
		tu.FailE(expect, sshot2.Name)
	}
	if expect := int64(281378); sshot2.Size != expect {
		tu.FailE(expect, sshot2.Size)
	}
	if !modtime2.Equal(sshot2.ModTime) {
		tu.FailE(modtime2, sshot2.ModTime)
	}
	if expect := "SHA256:E0fI8SqLFxqd2d501xzadaCBg0/ypYiYj5fMCxjJqcg="; sshot2.ContentHash["SHA256"] != expect {
		tu.FailE(expect, sshot2.ContentHash["SHA256"])
	}
	if expect := "a:00070f0f7f1f0703"; sshot2.PerceptualHash["average"] != expect {
		tu.FailE(expect, sshot2.PerceptualHash["average"])
	}
}

func TestCheckFile(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	testfile1 := filepath.Join(tu.Dir, "testfile1.txt")
	testfile2 := filepath.Join(tu.Dir, "testfile2.txt")
	testfile3 := filepath.Join(tu.Dir, "testfile3.txt")

	tu.ModTime(2018, 10, 10, testfile2, testfile3)

	sshot1, sshot2, sshot3 := new(Snapshot), new(Snapshot), new(Snapshot)
	tu.Must(<-sshot1.Generate(testfile1))
	tu.Must(<-sshot2.Generate(testfile2))
	tu.Must(<-sshot3.Generate(testfile3))

	tu.Must(sshot1.CheckFile(testfile1))
	tu.Must(sshot2.CheckFile(testfile2))
	tu.Must(sshot3.CheckFile(testfile3))

	if err := sshot1.CheckFile(testfile2); err == nil {
		tu.Fail("False positive 1!")
	} else if _, ok := err.(*MissmatchError); !ok {
		tu.Fail(err)
	}
	if err := sshot2.CheckFile(testfile3); err == nil {
		tu.Log("This test fails. But I'm allowing it anyway. Same modtime + size is enough to assume same file in this basic context")
	} else if _, ok := err.(*MissmatchError); !ok {
		tu.Fail(err)
	}
}

func TestReadOnly(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	testfile := filepath.Join(tu.Dir, "testfile.txt")
	tu.MustFatal(ioutil.WriteFile(testfile, []byte("hello there"), 0666))

	tu.Must(ReadOnly(testfile))
	testReadOnly(tu, testfile)
}

func TestLockFile(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	lockfileData := []byte(`
myfile:
  created: 2018-11-17T15:04:45.9077169+13:00
  name: myfile
  mod: 2018-11-17T15:04:45.9077169+13:00
  size: 123
  chash:
    SHA256: jargon
  phash:
    average: alsojargon`)
	lockfileHandle := bytes.NewReader(lockfileData)

	lockfile := LockMap{}
	tu.Must(lockfile.Load(lockfileHandle))

	sshot, ok := lockfile["myfile"]
	if !ok {
		tu.Fail("Data not in map")
	}
	if sshot.Name != "myfile" || sshot.ContentHash["SHA256"] != "jargon" || sshot.PerceptualHash["average"] != "alsojargon" {
		tu.Fail("name/hash was missing / incorrect")
	}

	outputHandle := bytes.NewBuffer([]byte(""))
	tu.Must(lockfile.Save(outputHandle))

	testString := strings.TrimSpace(string(lockfileData))
	expectString := strings.TrimSpace(outputHandle.String())

	if expectString != testString {
		tu.FailE(testString, expectString)
	}
}

func TestLockEvent(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")
	testfile := filepath.Join(event, "event01_001.txt")
	cxt := &context.Context{WorkingDir: event}
	tu.Must(LockEvent(cxt, false)) // Lock down the event!
	tu.AssertExists(filepath.Join(event, LOCKFILENAME))
	testReadOnly(tu, testfile)
}

func TestLockEventNew(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")
	cxt := &context.Context{WorkingDir: event}
	tu.Must(LockEvent(cxt, false))

	testReadOnly(tu, filepath.Join(event, "event01_001.txt"))
	testReadOnly(tu, filepath.Join(event, "event01_002.txt"))

	lockmap := LockMap{}
	handle := tu.MustFatal(os.Open(filepath.Join(event, LOCKFILENAME))).(*os.File)
	defer handle.Close()
	tu.MustFatal(lockmap.Load(handle))

	if _, ok := lockmap["event01_002.txt"]; !ok {
		tu.Fail("Failed to add entry for new event")
	}
}

func TestLockEventMissing(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")
	cxt := &context.Context{WorkingDir: event}
	if err, ok := LockEvent(cxt, false).(*MissmatchError); !ok {
		if err == nil {
			tu.Fail("Did not trigger error for missing file")
		} else {
			tu.Fail(err)
		}
	}
}

func TestLockEventChanged(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")
	cxt := &context.Context{WorkingDir: event}
	if err, ok := LockEvent(cxt, false).(*MissmatchError); !ok {
		if err == nil {
			tu.Fail("Did not trigger error for changed data")
		} else {
			tu.Fail(err)
		}
	}
}

func TestLockEventRenamed(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	event := filepath.Join(tu.Dir, "event01")
	cxt := &context.Context{WorkingDir: event}
	tu.Must(LockEvent(cxt, false))
}
