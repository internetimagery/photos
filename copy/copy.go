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
	sourceInfo, err := os.Stat(sourceDir) // Source must exist and be a directory
	if err != nil {
		return err
	}
	if !sourceInfo.IsDir() {
		return fmt.Errorf("Source is not a directory '%s'", sourceDir)
	}
	destInfo, err := os.Stat(destinationDir) // Destination will be created if doesn't exist. But cannot be a file
	if err == nil && !destInfo.IsDir() {
		return fmt.Errorf("destination is not a directory")
	} else if !os.IsNotExist(err) {
		return err
	}
	var tempRoot string
	if err == nil { // Destination already exists, make tempfile there
		tempRoot = filepath.Dir(destinationDir)
	} else { // Fall back to using source location
		tempRoot = filepath.Dir(sourceDir)
	}

	// Create temporary working folder to copy into initially
	tmpDir, err := ioutil.TempDir(filepath.Dir(tempRoot), "tempcopytree")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Run through all files and kick off copies
	if err = filepath.Walk(sourceDir, func(sourcePath string, info os.FileInfo, err error) error {
		if sourcePath == sourceDir {
			return nil // Ignore root. We know it exists already!
		}

		// Gather our source and destination file paths
		relPath, err := filepath.Rel(sourceDir, sourcePath)
		if err != nil {
			return err
		}
		destPath := filepath.Join(tmpDir, relPath)

		// Don't worry about parralellizing directory creation. Get that over with quickly in serial
		if info.IsDir() {
			if err = os.Mkdir(destPath, info.Mode().Perm()); err != nil {
				return err
			}
		} else {
			// TODO: Store these done channels. Check for them to finish later on
			// TODO: Consider putting in another channel that stops execution on error
			if err = <-File(sourcePath, destPath); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// Put everything into its right place!
	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(destinationDir, sourceInfo.Mode().Perm()); err != nil {
		return err
	}
	fmt.Println("FILES", files)
	for _, file := range files {
		sourcePath := filepath.Join(tmpDir, file.Name())
		destPath := filepath.Join(destinationDir, file.Name())
		fmt.Println(sourcePath, destPath)
		if err := os.Rename(sourcePath, destPath); err != nil {
			return err
		}
	}
	return nil
}
