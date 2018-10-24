package copy

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// File : Convenience wrapper for copyfile. Sets up connection channel between the two. Can be used in serial too
func File(source, destination string) chan error {
	done := make(chan error)
	go copyfile(source, destination, done)
	return done
}

// copyfile : Copy a file from one place to another. Retain permissions (though not owner)
func copyfile(source, destination string, done chan error) {
	var err error // Reuse error variable everywhere for return var
	defer func() {
		done <- err
	}()

	// Check our source exists, and our destination doesnt
	if _, err = os.Stat(source); err != nil {
		return
	}
	if _, err = os.Stat(destination); !os.IsNotExist(err) {
		if err == nil {
			err = fmt.Errorf("Destination file exists '%s'", destination)
		}
		return
	}

	// Open our sourcefile, and a temporary file in destination location
	sourceHandle, err := os.Open(source)
	if err != nil {
		return
	}
	defer sourceHandle.Close()

	destinationHandle, err := ioutil.TempFile(filepath.Dir(destination), "tempcopy")
	if err != nil {
		return
	}
	defer func(name string) { // Cleanup!
		if _, err := os.Stat(name); err != nil {
			os.Remove(name)
		}
	}(destinationHandle.Name())

	// Copy data across
	if _, err = io.Copy(destinationHandle, sourceHandle); err != nil {
		destinationHandle.Close()
		return
	}
	if err = destinationHandle.Close(); err != nil {
		return
	}

	// Set permissions and modification time
	const perm = 0644 // // TODO: query permissions to get this value
	if err = os.Chmod(destinationHandle.Name(), perm); err != nil {
		return
	}

	// Finally, set destination to its final resting place!
	err = os.Rename(destinationHandle.Name(), destination)
}

// 	const perm = 0644
// 	if err := os.Chmod(tmp.Name(), perm); err != nil {
// 		os.Remove(tmp.Name())
// 		return err
// 	}
// 	if err := os.Rename(tmp.Name(), dst); err != nil {
// 		os.Remove(tmp.Name())
// 		return err
// 	}
// 	return nil
// }
