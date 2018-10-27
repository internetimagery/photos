package copy

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// createDummy : Create a dummy file as a placeholder for a future copy
func createDummyFile(path string) error {
	modTime := time.Now().Add(time.Hour * 24 * 365) // Dummy file flagged by having been modified in the future
	handle, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	handle.Close()
	return os.Chtimes(path, modTime, modTime)
}

// createDummyDir : Create a dummy directory
func createDummyDir(path string) error {
	modTime := time.Now().Add(time.Hour * 24 * 365) // Dummy file flagged by having been modified in the future
	if err := os.Mkdir(path, 0700); err != nil {
		return err
	}
	return os.Chtimes(path, modTime, modTime)
}

// isDummy : The counterpart to createDummy. Checks if a given file is considered a dummy
func isDummy(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	// check dir
	if info.IsDir() && time.Now().Add(time.Hour*24*7).Before(info.ModTime()) {
		if files, err := ioutil.ReadDir(path); err == nil && len(files) == 0 {
			return true
		}
	} else if info.Mode().IsRegular() && info.Size() == 0 && time.Now().Add(time.Hour*24*7).Before(info.ModTime()) {
		// check file
		return true
	}
	return false
}

// cleanDummy : Remove all dummy files in directory
func cleanDummy(path string) error {
	return filepath.Walk(path, func(fileName string, info os.FileInfo, err error) error {
		if isDummy(fileName) {
			return os.Remove(fileName)
		}
		return nil
	})
}

// File : Convenience wrapper for copyfile. Sets up connection channel between the two. Can be used in serial too
func File(source, destination string) chan error {
	done := make(chan error)
	go func() {
		var err error
		defer func() {
			done <- err
		}()

		// Check our source exists
		sourceInfo, err := os.Stat(source)
		if err != nil {
			return
		} else if !sourceInfo.Mode().IsRegular() {
			err = fmt.Errorf("Source not a regular file '%s'", source)
			return
		}

		// Lock our destination with a dummy file
		if err = createDummyFile(destination); err != nil {
			return
		}
		defer func(name string) {
			if isDummy(name) {
				os.Remove(name)
			}
		}(destination)

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

		// Finally, set destination to its final resting place, replacing the dummy file!
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

	// TODO: create walk function to determine which folder exists for cleanup
	// TODO: make destination up front
	// TODO: mock out directories with dummy files
	// TODO: don't just rename root level files. What happens if a folder already contains a file. It'll be squashed.

	// TODO: strategy:
	// TODO: validate source
	// TODO: validate dest
	// TODO: make dest if needed (then remove dest if return with error)
	// TODO: walk source
	// TODO: make dummy files in dest. Run cleanup after to remove stray dummies
	// TODO: initiate file copying

	// wait for copies to finish
	// check errors
	// walk tempfile, replace dummies with real deal

	// Validate our input exists and is a directory
	sourceInfo, err := os.Stat(sourceDir)
	if err != nil {
		return err
	}
	if !sourceInfo.IsDir() {
		return fmt.Errorf("Source is not a directory '%s'", sourceDir)
	}

	// Validate destination either exists, and is a directory
	// or does not exist, needs to be created and parent dir exists
	destInfo, err := os.Stat(destinationDir)
	if err == nil && !destInfo.IsDir() {
		return fmt.Errorf("destination is not a directory")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = os.Mkdir(destinationDir, 0700)
	if err == nil { // We created this directory, cleanup after
		defer func() {
			if err == nil { // No error? Set permissions of file to match source.
				err = os.Chmod(destinationDir, sourceInfo.Mode().Perm())
			} else { // Was an error. We created this directory. Clean it up.
				os.RemoveAll(destinationDir)
			}
		}()
	} else if !os.IsExist(err) { // Some other error
		return err
	}

	// Create temporary working folder to copy into initially
	tmpDir, err := ioutil.TempDir(destinationDir, "tempcopytree")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Run through all files. Prep dummy files, and kick off copies.
	copies := []chan error{}
	defer cleanDummy(sourceDir) // Ensure all dummyfiles are cleaned
	if err = filepath.Walk(sourceDir, func(sourcePath string, info os.FileInfo, err error) error {
		if sourcePath == sourceDir {
			return nil // Ignore root. We know it exists already, thanks!
		}

		// Gather our source and destination file paths
		relPath, err := filepath.Rel(sourceDir, sourcePath)
		if err != nil {
			return err
		}
		destPath := filepath.Join(tmpDir, relPath)
		dummyPath := filepath.Join(destinationDir, relPath)

		// Don't worry about parralellizing directory creation. Get that over with quickly in serial
		if info.IsDir() {
			if err = createDummyDir(dummyPath); err != nil {
				return err
			}
			if err = os.Mkdir(destPath, info.Mode().Perm()); err != nil {
				return err
			}
		} else {
			if err = createDummyFile(dummyPath); err != nil {
				return err
			}
			// TODO: Consider putting in another channel that stops execution on error
			copies = append(copies, File(sourcePath, destPath))
		}
		return nil
	}); err != nil {
		return err
	}

	// Ensure all copies have finished.
	for _, done := range copies {
		if err = <-done; err != nil {
			return err
		}
	}

	// Put everything into its right place!
	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(destinationDir, sourceInfo.Mode().Perm()); err != nil {
		return err
	}
	for _, file := range files {
		sourcePath := filepath.Join(tmpDir, file.Name())
		destPath := filepath.Join(destinationDir, file.Name())
		if err := os.Rename(sourcePath, destPath); err != nil {
			return err
		}
	}
	return nil
}
