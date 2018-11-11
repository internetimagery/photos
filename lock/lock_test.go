package lock

import (
	"bytes"
	"testing"

	"github.com/internetimagery/photos/testutil"
)

func TestGenerateContentHash(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	buff := bytes.NewReader([]byte("SOME DATA HERE"))
	expectHash := "j7BgIUq2w472YYetmry+ieE0D3kqaVRdU6Ri6uq2hTY=" // MD5

	testHash, err := GenerateContentHash("MD5", buff)
	if err != nil {
		tu.Fail(err)
	}
	if expectHash != testHash {
		tu.FailE(expectHash, testHash)
	}
}
