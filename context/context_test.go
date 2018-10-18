package context

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/testutil"
)

func TestNewContext(t *testing.T) {
	tmpDir := testutil.NewTempDir(t, "TestNewContext")
	defer tmpDir.Close()

	if _, err := NewContext(tmpDir.Dir); !os.IsNotExist(err) {
		fmt.Println(err)
		t.Fail()
	}
}

func TestContext(t *testing.T) {
	configData := []byte(`{
	"compress": [
		["*", "some command"]
	]
}`)

	tmpDir := testutil.NewTempDir(t, "TestContext")
	defer tmpDir.Close()

	// Make a couple files
	workingDir := filepath.Join(tmpDir.Dir, "some-event")
	rootConf := filepath.Join(tmpDir.Dir, ROOTCONF)

	err := os.Mkdir(workingDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(rootConf, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Start within event directory
	cxt, err := NewContext(workingDir)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Check we reached root file
	if cxt.Root != tmpDir.Dir {
		fmt.Println("Couldn't find config file.", cxt.Root)
		t.Fail()
	}
}

func TestContextEnv(t *testing.T) {
	tmpDir := testutil.NewTempDir(t, "TestContextEnv")
	defer tmpDir.Close()

	err := ioutil.WriteFile(filepath.Join(tmpDir.Dir, ROOTCONF), []byte("{}"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Set environment var
	os.Setenv("TESTENV", "SUCCESS")

	// Build context and check environment var came through
	cxt, err := NewContext(tmpDir.Dir)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if cxt.Env["TESTENV"] != "SUCCESS" {
		fmt.Println("Env was not passed into context")
		t.Fail()
	}
}

func TestContextPrepCommand(t *testing.T) {
	cxt := &Context{Env: map[string]string{
		"TESTENV": "VALUE",
	}}

	expectCommand := "echo $TESTENV"
	command, err := cxt.PrepCommand(expectCommand)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(command.Args) != 2 && command.Args[1] != "VALUE" {
		fmt.Println("Got args", command.Args)
		t.Fail()
	}

}
