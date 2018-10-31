// Adding simple copy functionality to windows for testing
package main

import (
	"fmt"
	"os"
	"runtime"
)

func makeFile(path string) error {
	// Create our file
	handle, err := os.OpenFile(path, os.O_WRONLY|os.O_EXCL|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return handle.Close()
}

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println("Sorry this tool is meant to be run on windows for testing")
		return
	}

	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Expecting: exec.exe <source>")
	}

	source := os.Args[1]

	if err := makeFile(source); err == nil {
		os.Exit(0)
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}
