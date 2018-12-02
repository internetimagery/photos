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
	"strings"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/internetimagery/photos/format"
	yaml "gopkg.in/yaml.v2"
)

// LOCKFILENAME : Name of file displaying the locked state of an event
const LOCKFILENAME = "locked"

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

// CheckFile : Check if two snapshots refer to the same file
func (sshot *Snapshot) Equals(other *Snapshot) error {

	if other.Size != sshot.Size {
		return &MissmatchError{"Snapshot size does not match: " + other.Name}
	}
	if other.ModTime == sshot.ModTime { // There needs to be some margin of error here. But how much? What does rclone do?
		// Roughly conclude a match!
		return nil
	}
	for hashType, hashVal := range sshot.ContentHash {
		if otherHashType, ok := other.ContentHash[hashType]; ok { // Compare similar hashes!
			if hashVal != otherHashType {
				return &MissmatchError{fmt.Sprintf("Snapshot hash does not match: '%s' %s", hashType, otherHashType)}
			}
		} else {
			continue
		}
	}
	return &MissmatchError{"Snapshots do not have comparable hashes: " + other.Name}
}

// Save : Save snapshot data!
func (sshot *Snapshot) Save(handle io.Writer) error {
	data, err := yaml.Marshal(sshot)
	if err != nil {
		return err
	}
	_, err = handle.Write(data)
	return err
}

// Load : Load snapshot data
func (sshot *Snapshot) Load(handle io.Reader) error {
	data, err := ioutil.ReadAll(handle)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &sshot)
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

// LoadSnapshot : load a snapshot file. Fixing up Name variable if missing.
func LoadSnapshot(filename string) (*Snapshot, error) {
	sshot := &Snapshot{}
	ext := filepath.Ext(filename)
	if ext != "yaml" {
		return sshot, fmt.Errorf("File missing correct extension '%s'", filename)
	}
	handle, err := os.Open(filename)
	if err != nil {
		return sshot, err
	}
	if err = sshot.Load(handle); err != nil {
		return sshot, err
	}
	if strings.TrimSpace(sshot.Name) == "" {
		sourceBase := filepath.Base(filename)
		sshot.Name = sourceBase[len(sourceBase)-len(ext):]
	}
	return sshot, nil
}

// SnapshotManager : Map snapshots to their corresponding files
type SnapshotManager map[string]*Snapshot

// NewSnapshotManager : Load up an existing or new mapping of snapshot files
func NewSnapshotManager(eventname string) (SnapshotManager, error) {
	newMap := SnapshotManager{}
	lockname := filepath.Join(eventname, LOCKFILENAME)
	lockfiles, err := ioutil.ReadDir(lockname)
	if os.IsNotExist(err) {
		return newMap, nil
	} else if err != nil {
		return newMap, err
	}

	// Run through existing lock files and map them
	for _, lockfile := range lockfiles {
		if err = newMap.LoadSnapshot(filepath.Join(eventname, lockfile.Name())); err != nil {
			return newMap, err
		}
	}
	return newMap, nil
}

// LoadSnapshot : load a snapshot file and map it
func (snapmap SnapshotManager) LoadSnapshot(filename string) error {
	ext := filepath.Ext(filename)
	if ext != "yaml" {
		return fmt.Errorf("File missing correct extension '%s'", filename)
	}
	handle, err := os.Open(filename)
	if err != nil {
		return err
	}
	sshot := &Snapshot{}
	if err = sshot.Load(handle); err != nil {
		return err
	}
	sourceFile := sshot.Name
	if sourceFile == "" {
		sourceBase := filepath.Base(filename)
		sourceFile = sourceBase[len(sourceBase)-len(ext):]
	}
	snapmap[sourceFile] = sshot
	return nil
}

// AddSnapshot : Add a new snapshot
func (snapmap SnapshotManager) AddSnapshot(filename string) chan error {
	sshot := &Snapshot{}

}

// LockEvent : Attempt to lock event. If lock exists, check for any changes and update lock. If force, version up file and lock
func LockEvent(directoryname string, force bool) error {
	// Grab media from within file
	mediaList, err := format.NewEvent(directoryname).GetMedia()
	if err != nil {
		return err
	}
	if len(mediaList) == 0 { // Nothing to do!
		return nil
	}

	// Check path to lock folder exists, and create if it doesn't
	lockPath := filepath.Join(directoryname, LOCKFILENAME)
	if info, err := os.Stat(lockPath); os.IsNotExist(err) {
		if err = os.Mkdir(lockPath, 0755); err != nil {
			return err
		}
	} else if !info.IsDir() {
		return fmt.Errorf("file exists with the same name '%s', please remove it", LOCKFILENAME)
	} else if err != nil {
		return err
	}

	// Load a mapping of our lock files
	lockmap, err := NewSnapshotMap(directoryname)
	if err != nil {
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
				return err // If not a missmatch error, then it's some other problem
			} else if force {
				// TODO: ensure old data sticks around in case there is a backup somewhere
				// TODO: version up file, and check if no existing file is there. If so, version up again
				// TODO: Add new snapshot for this new file
			} else {
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

	// Refresh readonly on existing files
	for filename := range checkFiles {
		err := ReadOnly(filename)
		if err != nil {
			return err
		}
	}

	return nil
}
