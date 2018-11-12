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

// TODO: test different image types (eg png)
// TODO: test error on non-image type
func TestGeneratePercetualHash(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	handle := tu.MustFatal(os.Open(filepath.Join(tu.Dir, "testimg.jpg"))).(*os.File)
	defer handle.Close()

	expectHash := "a:00070f0f7f1f0703"
	if !tu.Must(IsSamePerceptualHash(expectHash, expectHash)).(bool) {
		tu.Fail("Well hash comparison failed...")
	}
	if tu.Must(IsSamePerceptualHash(expectHash, "a:00070f0f7f1f1233")).(bool) {
		tu.Fail("False positive hash comparison")
	}
	if _, err := IsSamePerceptualHash("am I a hash?", expectHash); err == nil {
		tu.Fail("Succeeded on bad first argument")
	}
	if _, err := IsSamePerceptualHash(expectHash, "am I a hash?"); err == nil {
		tu.Fail("Succeeded on bad second argument")
	}

	testHash := tu.Must(GeneratePerceptualHash("average", handle)).(string)
	if !tu.Must(IsSamePerceptualHash(expectHash, testHash)).(bool) {
		tu.FailE(expectHash, testHash)
	}
}
