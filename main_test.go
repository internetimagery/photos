package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/testutil"
)

func TestQuestion(t *testing.T) {
	defer testutil.UserInput(t, "y\n")()
	if !question() {
		fmt.Println("Question did not pass with 'y'")
		t.Fail()
	}

	defer testutil.UserInput(t, "n\n")
	if question() {
		fmt.Println("Question passed with 'n'")
		t.Fail()
	}

}

// Test init on a clean directory
func TestInitClean(t *testing.T) {
	tmpDir := testutil.NewTempDir(t, "TestInitClean")
	defer tmpDir.Close()

	defer testutil.UserInput(t, "y\n")()

	if err := run(tmpDir.Dir, []string{"exe", "init", "projectname"}); err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if _, err := os.Stat(filepath.Join(tmpDir.Dir, context.ROOTCONF)); err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// err := os.Chdir(tmpDir.Dir)
	// if err != nil {
	// 	fmt.Println(err)
	// 	t.Fatal(err)
	// }

	// fmt.Println(os.Args)

	// Generate and run command
	// executable := generateExec(tmpDir, t)
	// fmt.Println(os.Stat(executable))
	// cmd := exec.Command(os.Args[0], "init")
	// cmd.Dir = tmpDir
	// cmd.Stdout = os.Stdout
	// cmd.Stdin = strings.NewReader("n\n")
	// output, err := cmd.CombinedOutput()
	// err = cmd.Run()
	//
	// if err != nil {
	// 	// fmt.Println(string(output))
	// 	fmt.Println(err)
	// 	t.Fail()
	// }

	// fmt.Println(os.Stat(executable))
	//
	// // Test init in empty dir
	// com := exec.Command(executable, "init")
	// com.Stdin = bytes.NewReader([]byte("y"))
	// output, err := com.Output()
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Println("Running")
	// 	t.Fail()
	// }
	// fmt.Println(output)
}
