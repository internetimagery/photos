package lock

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/internetimagery/photos/context"
	"github.com/internetimagery/photos/format"
	yaml "gopkg.in/yaml.v2"
)

// TODO: generate phash
// TODO: make lock object to contain information
// TODO: impliment checking
// TODO: impliment serializing lock data
// TODO: make function to set files readonly, linux/osx/windows

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
func GeneratePerceptualHash(hashType string, handle io.Reader) (string, error) {
	img, err := jpeg.Decode(handle)
	if err != nil {
		return "", err
	}
	switch hashType {
	case "average":
		hash, err := goimagehash.AverageHash(img)
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
	return dist < 2, nil
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
		} else if _, ok := err.(jpeg.FormatError); ok {
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

// LockFile : Format and usage of locked file
type LockFile struct {
	Snapshots map[string]*Snapshot // Basic representations of the files
}

// Save : Save lockfile data!
func (lock *LockFile) Save(handle io.Writer) error {
	data, err := yaml.Marshal(lock)
	if err != nil {
		return err
	}
	_, err = handle.Write(data)
	return err
}

// Load : Load lockfile data
func (lock *LockFile) Load(handle io.Reader) error {
	data, err := ioutil.ReadAll(handle)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, lock)
}

// LockEvent : Attempt to lock event. If lock exists, check for any changes and update lock.
func LockEvent(cxt *context.Context, force bool) error {
	// Grab media from within file
	mediaList, err := format.GetMediaFromDirectory(cxt.WorkingDir)
	if err != nil {
		return err
	}
	// If there is nothing to use, ignore!
	if len(mediaList) == 0 {
		return nil
	}

	// Load lockfile snapshot data
	lockfilePath := filepath.Join(cxt.WorkingDir, LOCKFILENAME)
	lockfile := LockFile{}
	if handle, err := os.Open(lockfilePath); err == nil {
		err = lockfile.Load(handle)
		handle.Close()
		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	// Verify all files. Kick off new snapshots
	jobs := []chan error{}
	for _, media := range mediaList {
		name := filepath.Base(media.Path)
		sshot, ok := lockfile.Snapshots[name]
		if ok && !force { // File exists, compare snapshot, unless we are force locking things
			if err = sshot.CheckFile(media.Path); err != nil {
				return err
			}
		} else if media.Index > 0 { // Ignore unformatted files
			sshot = new(Snapshot)
			lockfile.Snapshots[name] = sshot
			jobs = append(jobs, sshot.Generate(media.Path))
		}
	}

	if len(jobs) > 0 { // We have added some things to the lockfile!
		// Finish up jobs
		for _, job := range jobs {
			if err = <-job; err != nil {
				return err
			}
		}

		// Make files readonly
		for file := range lockfile.Snapshots {
			err := ReadOnly(filepath.Join(cxt.WorkingDir, file))
			if err != nil {
				return err
			}
		}

		// Save lockfile!
		handle, err := os.OpenFile(lockfilePath, os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer handle.Close()
		if err = lockfile.Save(handle); err != nil {
			return err
		}
	}

	return nil
}
