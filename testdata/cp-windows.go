// Adding simple copy functionality to windows for testing
package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

func copyfile(source, destination string) error {
	// Open our sourcefile
	sourceHandle, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceHandle.Close()
	sourceInfo, err := sourceHandle.Stat()
	if err != nil {
		return err
	}

	// Create our desination
	destinationHandle, err := os.OpenFile(destination, os.O_WRONLY|os.O_EXCL|os.O_CREATE, sourceInfo.Mode().Perm())
	if err != nil {
		return err
	}
	defer func() { // Remove file if error
		if err != nil {
			os.Remove(destination) // Error handled? Nope...
		}
	}()

	// Copy data across
	if _, err = io.Copy(destinationHandle, sourceHandle); err != nil {
		destinationHandle.Close() // Error not handled, argh
		return err
	}
	if err = destinationHandle.Close(); err != nil {
		return err
	}

	// Set modification time
	if err = os.Chtimes(destinationHandle.Name(), sourceInfo.ModTime(), sourceInfo.ModTime()); err != nil {
		return err
	}
	return nil
}

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println("Sorry this tool is meant to be run on windows for testing")
		return
	}

	if len(os.Args) != 3 {
		fmt.Println("Wrong number of arguments. Expecting: exec.exe <source> <dest>")
	}

	source := os.Args[1]
	destination := os.Args[2]

	if err := copyfile(source, destination); err == nil {
		os.Exit(0)
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}
