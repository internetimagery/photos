package copy

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// File : Convenience wrapper for copyfile. Sets up connection channel between the two. Can be used in serial too
func File(source, destination string) error {

	// Check our source exists, and our destination doesnt
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	} else if !sourceInfo.Mode().IsRegular() {
		err = fmt.Errorf("Source not a regular file '%s'", source)
		return err
	}
	if _, err = os.Stat(destination); !os.IsNotExist(err) {
		if err == nil {
			err = fmt.Errorf("Destination file exists '%s'", destination)
		}
		return err
	}

	// Open our sourcefile, and a temporary file in destination location
	sourceHandle, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceHandle.Close()

	destinationHandle, err := ioutil.TempFile(filepath.Dir(destination), "tempcopy")
	if err != nil {
		return err
	}
	defer func(name string) { // Cleanup!
		if _, err := os.Stat(name); err != nil {
			os.Remove(name)
		}
	}(destinationHandle.Name())

	// Copy data across
	if _, err = io.Copy(destinationHandle, sourceHandle); err != nil {
		destinationHandle.Close()
		return err
	}
	if err = destinationHandle.Close(); err != nil {
		return err
	}

	// Set permissions and modification time
	if err = os.Chtimes(destinationHandle.Name(), sourceInfo.ModTime(), sourceInfo.ModTime()); err != nil {
		return err
	}
	perm := sourceInfo.Mode().Perm()
	if err = os.Chmod(destinationHandle.Name(), perm); err != nil {
		return err
	}

	// Finally, set destination to its final resting place!
	return os.Rename(destinationHandle.Name(), destination)

	// // Last minute permissions change if on windows
	// if runtime.GOOS == "windows" {
	// 	err = acl.Chmod(destinationHandle.Name(), perm)
	// }
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
		} else {
			return err
		}
	}

	poolSize := 5

	jobs := make(chan [2]string, 100)
	errors := make(chan string, 100)
	doneList := []chan error{}
	errorList := []string{}

	for i := 0; i < poolSize; i++ {
		done := make(chan error)
		go treeworker(jobs, done)
		doneList = append(doneList, done)
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
				errorList = append(errorList, err.Error())
			}
			return nil
		}

		// Add job
		jobs <- [2]string{sourcePath, destPath}
		return nil
	})

	// Collect our errors (hopefully none!)
	for _, done := range doneList {
		err = <-done
		if err != nil {
			errorList = append(errorList, err.Error())
		}
	}
	for i := len(errors); i > 0; i = len(errors) {
		errorList = append(errorList, <-errors)
	}

	if len(errorList) != 0 {
		return fmt.Errorf(strings.Join(errorList, "\n"))
	} else {
		return nil
	}
}

// treeworker : Run jobs on behalf of the Tree caller
func treeworker(jobs chan [2]string, done chan error) {
	errList := []string{}
	for job, ok := <-jobs; ok; job, ok = <-jobs {
		log.Printf("Copying '%s' --> '%s'", job[0], job[1])
		if err := File(job[0], job[1]); err != nil {
			errList = append(errList, err.Error())
		}
	}
	if len(errList) != 0 {
		done <- fmt.Errorf(strings.Join(errList, "\n"))
	} else {
		done <- nil
	}
}
