package lock

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/testutil"
)

func TestGenerateContentHash(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	buff := bytes.NewReader([]byte("SOME DATA HERE"))
	expectHash := "SHA256:j7BgIUq2w472YYetmry+ieE0D3kqaVRdU6Ri6uq2hTY=" // MD5

	testHash := tu.Must(GenerateContentHash("SHA256", buff)).(string)
	if expectHash != testHash {
		tu.FailE(expectHash, testHash)
	}
}

func testgeneratepercetualhash(t *testing.T) {
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

	testHash1 := tu.Must(GeneratePerceptualHash("average", handle1)).(string)
	testHash2 := tu.Must(GeneratePerceptualHash("average", handle2)).(string)
	testHash3 := tu.Must(GeneratePerceptualHash("average", handle3)).(string)
	testHash4 := tu.Must(GeneratePerceptualHash("average", handle4)).(string)
	_, err := GeneratePerceptualHash("average", handle5)
	if err == nil {
		tu.Fail("Allowed unsupported img")
	}
	_, err = GeneratePerceptualHash("average", handle6)
	if err == nil {
		tu.Fail("Allowed unsupported file")
	}

	if !tu.Must(IsSamePerceptualHash(testHash1, testHash2)).(bool) ||
		!tu.Must(IsSamePerceptualHash(testHash1, testHash3)).(bool) ||
		!tu.Must(IsSamePerceptualHash(testHash2, testHash3)).(bool) {
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

	if sshot1.Name != "testimg1.txt" ||
		sshot1.Size != 24 ||
		!modtime1.Equal(sshot1.ModTime) ||
		sshot1.ContentHash["SHA256"] != "SHA256:h13POS/MwQ0SHVmJOSHgeN7+fM9ymIJZvdZt3nnLAqY=" ||
		sshot1.PerceptualHash != nil {
		tu.Fail("Invalid snapshot 1")
	}

	if sshot2.Name != "testimg2.jpg" ||
		sshot2.Size != 281378 ||
		!modtime2.Equal(sshot2.ModTime) ||
		sshot2.ContentHash["SHA256"] != "SHA256:E0fI8SqLFxqd2d501xzadaCBg0/ypYiYj5fMCxjJqcg=" ||
		sshot2.PerceptualHash["average"] != "a:00070f0f7f1f0703" {
		tu.Fail("Invalid snapshot 2")
	}
}
