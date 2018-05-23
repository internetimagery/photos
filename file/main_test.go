// safe file operations.

package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteOk(t *testing.T) {
	dir, err := ioutil.TempDir("", "file_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	data := "Here is some info"
	path := filepath.Join(dir, "ok.file")
	err := file.Write([]byte(data), path)
	if err != nil {
		t.Error(err)
	}
}
