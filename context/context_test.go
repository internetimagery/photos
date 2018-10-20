package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/testutil"
)

func TestNewContext(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.TempDir("TestNewContext")()

	if _, err := NewContext(tu.Dir); !os.IsNotExist(err) {
		tu.Fail(err)
	}
}

func TestContext(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.TempDir("TestContext")()

	configData := `{
	"compress": [
		["*", "some command"]
	]
}`

	// Make a couple files
	workingDir := filepath.Join(tu.Dir, "some-event")
	rootConf := filepath.Join(tu.Dir, ROOTCONF)

	if err := os.Mkdir(workingDir, 0755); err != nil {
		tu.Fatal(err)
	}

	tu.NewFile(rootConf, configData)

	// Start within event directory
	cxt, err := NewContext(workingDir)
	if err != nil {
		tu.Fail(err)
	}

	// Check we reached root file
	if cxt.Root != tu.Dir {
		tu.Fail("Couldn't find config file.", cxt.Root)
	}
}

func TestContextEnv(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.TempDir("TestContextEnv")()

	tu.NewFile(filepath.Join(tu.Dir, ROOTCONF), "{}")

	// Set environment var
	os.Setenv("TESTENV", "SUCCESS")

	// Build context and check environment var came through
	cxt, err := NewContext(tu.Dir)
	if err != nil {
		tu.Fail(err)
	}
	if cxt.Env["TESTENV"] != "SUCCESS" {
		tu.Fail("Env was not passed into context")
	}
}

func TestContextPrepCommand(t *testing.T) {
	tu := testutil.NewTestUtil(t)

	cxt := &Context{Env: map[string]string{
		"TESTENV": "VALUE",
	}}

	expectCommand := "echo $TESTENV"
	command, err := cxt.PrepCommand(expectCommand)
	if err != nil {
		tu.Fail(err)
	}
	if len(command.Args) != 2 && command.Args[1] != "VALUE" {
		tu.Fail("Got args", command.Args)
	}

}
