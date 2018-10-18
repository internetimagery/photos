package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/internetimagery/photos/testutil"
)

func generateExec(tmpDir string, t *testing.T) string {
	executable := filepath.Join(tmpDir, "main")
	com := exec.Command("go", "build", "-o", executable)
	output, err := com.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		t.Fatal(err)
	}
	return executable
}

// Test init on a clean directory
func TestInitClean(t *testing.T) {
	tmpDir := testutil.NewTempDir(t, "TestInitClean")
	defer tmpDir.Close()

	// err := os.Chdir(tmpDir.Dir)
	// if err != nil {
	// 	fmt.Println(err)
	// 	t.Fatal(err)
	// }

	fmt.Println(os.Args)

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
