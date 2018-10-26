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
	done := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			done <- err
		}()

		// Check our source exists, and our destination doesnt
		sourceInfo, err := os.Stat(source)
		if err != nil {
			return
		} else if !sourceInfo.Mode().IsRegular() {
			err = fmt.Errorf("Source not a regular file '%s'", source)
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
			if err = os.Remove(name); os.IsNotExist(err) {
				err = nil
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
		if err = os.Chtimes(destinationHandle.Name(), sourceInfo.ModTime(), sourceInfo.ModTime()); err != nil {
			return
		}
		perm := sourceInfo.Mode().Perm()
		if err = os.Chmod(destinationHandle.Name(), perm); err != nil {
			return
		}

		// Finally, set destination to its final resting place!
		err = os.Rename(destinationHandle.Name(), destination)

		// // Last minute permissions change if on windows
		// if runtime.GOOS == "windows" {
		// 	err = acl.Chmod(destinationHandle.Name(), perm)
		// }
	}()
	return done
}

// Tree : Copy files and directories recursively
func Tree(sourceDir, destinationDir string) error {

	// Validate our input
	info, err := os.Stat(sourceDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("Source is not a directory '%s'", sourceDir)
	}
	if _, err = os.Stat(destinationDir); !os.IsNotExist(err) {
		if err != nil {
			return fmt.Errorf("Destination exists already '%s'", destinationDir)
		}
		return err
	}

	// Run through all files and kick off copies
	filepath.Walk(sourceDir, func(sourcePath string, info os.FileInfo, err error) error {

		// Gather our source and destination file paths
		relPath, err := filepath.Rel(sourceDir, sourcePath)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destinationDir, relPath)

		// Don't worry about parralellizing directory creation. Get that over with quickly in serial
		if info.IsDir() {
			if err = os.Mkdir(destPath, info.Mode().Perm()); err != nil {
				return err
			}
		} else {
			// TODO: Store these done channels. Check for them to finish later on
			// TODO: Consider putting in another channel that stops execution on error
			// TODO: Put files into a temporary directory first, then move them over afterwards
			if err = <-File(sourcePath, destPath); err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}
