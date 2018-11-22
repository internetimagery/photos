package lock

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/internetimagery/photos/format"
	yaml "gopkg.in/yaml.v2"
)

// LOCKFILENAME : Name of file displaying the locked state of an event
const LOCKFILENAME = "locked.yaml"

// GenerateContentHash : Generate hash from content to compare contents
func GenerateContentHash(hashType string, handle io.Reader) (string, error) {
	switch hashType {
	case "SHA256":
		hasher := sha256.New()
		if _, err := io.Copy(hasher, handle); err != nil {
			return "", err
		}
		return hashType + ":" + base64.StdEncoding.EncodeToString(hasher.Sum([]byte{})), nil
	}
	return "", fmt.Errorf("Unknown hash format '%s'", hashType)
}

// GeneratePerceptualHash : Generate hash representing visual to compare imagery
func GeneratePerceptualHash(hashType string, handle io.ReadSeeker) (string, error) {
	img, err := jpeg.Decode(handle)
	if _, ok := err.(jpeg.FormatError); ok { // Not a jpeg
		handle.Seek(0, 0)
		img, err = png.Decode(handle)
		if err != nil { // Not png either...
			return "", image.ErrFormat
		}
	} else if err != nil {
		return "", image.ErrFormat
	}
	switch hashType {
	case "average":
		hash, err := goimagehash.AverageHash(img)
		if err != nil {
			return "", err
		}
		return hash.ToString(), nil
	case "difference":
		hash, err := goimagehash.DifferenceHash(img)
		if err != nil {
			return "", err
		}
		return hash.ToString(), nil
	case "perception":
		hash, err := goimagehash.PerceptionHash(img)
		if err != nil {
			return "", err
		}
		return hash.ToString(), nil
	}
	return "", fmt.Errorf("Unknown hash format '%s'", hashType)
}

// TODO: add similar images to test this on. Also differently compressed versions of same image.

// IsSamePerceptualHash : Hash comparison looking for equality
func IsSamePerceptualHash(hash1, hash2 string) (bool, error) {
	test1, err := goimagehash.ImageHashFromString(hash1)
	if err != nil {
		return false, err
	}
	test2, err := goimagehash.ImageHashFromString(hash2)
	if err != nil {
		return false, err
	}

	dist, err := test1.Distance(test2)
	if err != nil {
		return false, err
	}
	return dist < 5, nil
}

// Snapshot : Hold information about a particular files information
type Snapshot struct {
	Created        time.Time         `yaml:"created"` // Time this snapshot was created
	Name           string            `yaml:"name"`    // Base of path. Ie /one/two.three = two.three
	ModTime        time.Time         `yaml:"mod"`     // Modification time
	Size           int64             `yaml:"size"`    // Filesize!
	ContentHash    map[string]string `yaml:"chash"`   // Hash of the content
	PerceptualHash map[string]string `yaml:"phash"`   // Hash of the image
}

// Generate : Generate new snapshot data from file, with all the trimmings. "err <-&Snapshot{}.Generate(name)"
func (sshot *Snapshot) Generate(filename string) chan error {
	done := make(chan error)
	go func() {
		var err error
		defer func() {
			done <- err
		}()

		// Get a handle on things!
		handle, err := os.Open(filename)
		if err != nil {
			return
		}
		defer handle.Close()

		// Collect basic info on file!
		info, err := handle.Stat()
		if err != nil {
			return
		}

		// Basic info
		sshot.Name = info.Name()
		sshot.ModTime = info.ModTime()
		sshot.Size = info.Size()

		chash, err := GenerateContentHash("SHA256", handle) // SHA256 hardcoded for now
		if err != nil {
			return
		}
		sshot.ContentHash = map[string]string{"SHA256": chash}

		// TODO: This error needs to be managed for files that cannot have a phash (non-images)
		handle.Seek(0, 0)
		phash, err := GeneratePerceptualHash("average", handle)
		if err == nil { // SHA256 hardcoded for now
			sshot.PerceptualHash = map[string]string{"average": phash}
		} else if err == image.ErrFormat {
			err = nil // Ignore format error
		}
		sshot.Created = time.Now()
	}()
	return done
}

// MissmatchError : Error type for missmatches
type MissmatchError struct {
	err string
}

func (err *MissmatchError) Error() string {
	return err.err
}

// CheckFile : Check if a snapshot matches corresponding file. Return missmatch error if not matching
func (sshot *Snapshot) CheckFile(filename string) error {
	// Get a handle
	handle, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer handle.Close()
	info, err := handle.Stat()
	if err != nil {
		return err
	}

	if info.Size() != sshot.Size {
		return &MissmatchError{"Size does not match: " + filename}
	}
	if info.ModTime() == sshot.ModTime { // There needs to be some margin of error here. But how much? What does rclone do?
		// Roughly conclude a match!
		return nil
	}
	hash, err := GenerateContentHash("SHA256", handle) // SHA256 hardcoding for now
	if err != nil {
		return err
	}
	if hash != sshot.ContentHash["SHA256"] {
		return &MissmatchError{"Content does not match: " + filename}
	}
	return nil
}

// ReadOnly : Make file readonly
func ReadOnly(filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if err := os.Chmod(filename, info.Mode().Perm()&0444); err != nil {
		return err
	}
	return nil
}

// LockMap : Format and usage of locked file
type LockMap map[string]*Snapshot // Basic representations of the files

// Save : Save lockfile data!
func (lock *LockMap) Save(handle io.Writer) error {
	data, err := yaml.Marshal(lock)
	if err != nil {
		return err
	}
	_, err = handle.Write(data)
	return err
}

// Load : Load lockfile data
func (lock *LockMap) Load(handle io.Reader) error {
	data, err := ioutil.ReadAll(handle)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &lock)
}

// LockEvent : Attempt to lock event. If lock exists, check for any changes and update lock.
func LockEvent(directoryname string, force bool) error {
	// Grab media from within file
	mediaList, err := format.GetMediaFromDirectory(directoryname)
	if err != nil {
		return err
	}

	// Load lockfile snapshot data if it exists
	lockmapPath := filepath.Join(directoryname, LOCKFILENAME)
	lockmap := LockMap{}
	if handle, err := os.Open(lockmapPath); err == nil {
		err = lockmap.Load(handle)
		handle.Close()
		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	// Sort out our files!
	newFiles := map[string]struct{}{}
	checkFiles := map[string]*Snapshot{}
	removedFiles := map[string]struct{}{}
	for _, media := range mediaList { // Look through existing files
		if media.Index <= 0 { // Skip unformatted media
			continue
		}
		basename := filepath.Base(media.Path)
		if sshot, ok := lockmap[basename]; ok {
			checkFiles[media.Path] = sshot // File exists in folder and in locked list. Check for changes
		} else {
			newFiles[media.Path] = struct{}{} // New file not in locked list
		}
	}
	for basename := range lockmap {
		filename := filepath.Join(directoryname, basename)
		if _, ok := checkFiles[filename]; !ok {
			removedFiles[basename] = struct{}{} // File exists in list, but not in folder. It has been removed (or renamed)
		}
	}

	// If we have nothing to do... we're done!
	if len(newFiles) == 0 && len(checkFiles) == 0 && len(removedFiles) == 0 {
		return nil
	}

	// First, we'll make snapshots out of our new files
	newSnapshots := map[*Snapshot]chan error{}
	for filename := range newFiles {
		sshot := new(Snapshot)
		newSnapshots[sshot] = sshot.Generate(filename)
	}
	for sshot, job := range newSnapshots {
		if err = <-job; err != nil {
			return err
		}
		lockmap[sshot.Name] = sshot
	}

	// Next we'll check to see if any missing files are actually in the new snapshots (rename)
	for basename := range removedFiles {
		ok := false
		for sshot := range newSnapshots { // compare hashes (hard coded sha256 for now...)
			if lockmap[basename].ContentHash["SHA256"] == sshot.ContentHash["SHA256"] {
				ok = true // Looks like this file matches another new file. Transparently deal with the rename and continue
				delete(lockmap, basename)
			}
		}
		if !ok && !force {
			return &MissmatchError{"File was removed: " + basename}
		}
	}

	// Finally lets verify that our existing files are still ok!
	for filename, sshot := range checkFiles {
		if err = sshot.CheckFile(filename); err != nil {
			if _, ok := err.(*MissmatchError); !ok {
				return err
			} else if !force {
				return err
			}
		}
	}

	// Save lockmap!
	handle, err := os.OpenFile(lockmapPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer handle.Close()
	if err = lockmap.Save(handle); err != nil {
		return err
	}

	// Make new files readonly
	for filename := range newFiles {
		err := ReadOnly(filename)
		if err != nil {
			return err
		}
	}
	return nil
}
