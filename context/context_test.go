package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/testutil"
)

func TestNewContext(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	if _, err := NewContext(tu.Dir); !os.IsNotExist(err) {
		tu.Fail(err)
	}
}

func TestContext(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	defer tu.LoadTestdata()()

	workingDir := filepath.Join(tu.Dir, "some-event")
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
	defer tu.LoadTestdata()()

	// Set environment var
	os.Setenv("TESTENV", "SUCCESS")

	// Build context and check environment var came through
	cxt, err := NewContext(tu.Dir)
	if err != nil {
		tu.Fatal(err)
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
