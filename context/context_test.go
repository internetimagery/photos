package context

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewContext(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestNewContext")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	if _, err := NewContext(tmpDir); !os.IsNotExist(err) {
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

	tmpDir, err := ioutil.TempDir("", "TestContext")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Make a couple files
	workingDir := filepath.Join(tmpDir, "some-event")
	rootConf := filepath.Join(tmpDir, ROOTCONF)

	err = os.Mkdir(workingDir, 0755)
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
	if cxt.Root != tmpDir {
		fmt.Println("Couldn't find config file.", cxt.Root)
		t.Fail()
	}
}

func TestContextEnv(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestContextEnv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	err = ioutil.WriteFile(filepath.Join(tmpDir, ROOTCONF), []byte("{}"), 644)
	if err != nil {
		t.Fatal(err)
	}

	// Set environment var
	os.Setenv("TESTENV", "SUCCESS")

	// Build context and check environment var came through
	cxt, err := NewContext(tmpDir)
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
