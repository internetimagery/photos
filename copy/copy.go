package copy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// File : Copy a file. Can be used in serial or goroutine
func File(source, destination string) chan error {
	done := make(chan error)
	go func() {
		var err error
		defer func() {
			done <- err
		}()

		// Open our sourcefile
		sourceHandle, err := os.Open(source)
		if err != nil {
			return
		}
		defer sourceHandle.Close()
		sourceInfo, err := sourceHandle.Stat()
		if err != nil {
			return
		}

		// Create our desination
		destinationHandle, err := os.OpenFile(destination, os.O_WRONLY|os.O_EXCL|os.O_CREATE, sourceInfo.Mode().Perm())
		if err != nil {
			return
		}
		defer func() { // Remove file if error
			if err != nil {
				os.Remove(destination) // Error handled? Nope...
			}
		}()

		// Copy data across
		if _, err = io.Copy(destinationHandle, sourceHandle); err != nil {
			destinationHandle.Close() // Error not handled, argh
			return
		}
		if err = destinationHandle.Close(); err != nil {
			return
		}

		// Set modification time
		if err = os.Chtimes(destinationHandle.Name(), sourceInfo.ModTime(), sourceInfo.ModTime()); err != nil {
			return
		}

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

	// Map files to their channels
	type Media struct {
		Job        chan error
		SourcePath string
		SourceInfo os.FileInfo
		DestPath   string
		Err        error
	}
	jobs := []*Media{}

	// Run through all files and kick off copies.
	err = filepath.Walk(sourceDir, func(sourcePath string, info os.FileInfo, err error) error {
		if sourcePath == sourceDir {
			return nil // Ignore root. We know it exists already, thanks!
		}

		// Gather our source and destination file paths
		relPath, err := filepath.Rel(sourceDir, sourcePath)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destinationDir, relPath)

		// If we have a directory. Just make the damn thing! :P
		job := make(chan error, 1)
		if info.IsDir() {
			if err = os.Mkdir(destPath, 0755); err != nil {
				return err
			}
			job <- nil
		} else {
			job = File(sourcePath, destPath)
		}
		jobs = append(jobs, &Media{
			Job:        job,
			SourceInfo: info,
			SourcePath: sourcePath,
			DestPath:   destPath})
		return nil
	})
	defer func() { // If we had an error, clean up in reverse order
		if err != nil {
			for i := len(jobs) - 1; i >= 0; i-- {
				if jobs[i].Err == nil { // Don't clean up existing items. Only successful copies are required
					os.Remove(jobs[i].DestPath) // Cleaning up
				}
			}
		}
	}()

	// Ensure all in-progress copies have finished. Retain first error we come across.
	for _, media := range jobs {
		media.Err = <-media.Job // Set error from job
		if media.Err != nil && err == nil {
			err = media.Err
		}
	}
	if err != nil {
		return err
	}

	// Finally set permissions and modification times to match the source
	for i := len(jobs) - 1; i >= 0; i-- { // We can assume all copies have nil error at this point...
		job := jobs[i]
		if err = os.Chtimes(job.DestPath, job.SourceInfo.ModTime(), job.SourceInfo.ModTime()); err != nil {
			return err
		}
		if err = os.Chmod(job.DestPath, job.SourceInfo.Mode().Perm()); err != nil {
			return err
		}
	}

	// Done!
	return nil
}
