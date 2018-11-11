package context

import (
	"os"
	"os/exec"
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
	cxt := tu.Must(NewContext(workingDir)).(*Context)

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
	cxt := tu.Must(NewContext(tu.Dir)).(*Context)
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
	command := tu.Must(cxt.PrepCommand(expectCommand)).(*exec.Cmd)
	if len(command.Args) != 2 && command.Args[1] != "VALUE" {
		tu.Fail("Got args", command.Args)
	}
}

func TestContextAbsPath(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	cwd, _ := os.Getwd()
	cxt := &Context{WorkingDir: cwd}

	if expect, _ := filepath.Abs("/four/five"); expect != cxt.AbsPath("/four/five") {
		tu.FailE(expect, cxt.AbsPath("/four/five"))
	}
	if expect, _ := filepath.Abs("four/five"); expect != cxt.AbsPath("four/five") {
		tu.FailE(expect, cxt.AbsPath("four/five"))
	}
	if expect, _ := filepath.Abs("."); expect != cxt.AbsPath(".") {
		tu.FailE(expect, cxt.AbsPath("."))
	}
	if expect, _ := filepath.Abs("four/../../five"); expect != cxt.AbsPath("four/../../five") {
		tu.FailE(expect, cxt.AbsPath("four/../../five"))
	}
}
