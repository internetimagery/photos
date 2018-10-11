package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"
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

func TestInit(t *testing.T) {
	// TODO: Get this working...
	// wd, _ := os.Getwd()
	// tmpDir, err := ioutil.TempDir(wd, "phototest")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer os.RemoveAll(tmpDir)
	// executable := generateExec(tmpDir, t)
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
