package lock

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io"
	"time"
    "os"
    "path/filepath"
	"github.com/corona10/goimagehash"
)

// TODO: generate phash
// TODO: make lock object to contain information
// TODO: impliment checking
// TODO: impliment serializing lock data
// TODO: make function to set files readonly, linux/osx/windows

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
	Name           string            `yaml:"name"`  // Base of path. Ie /one/two.three = two.three
	ModTime        time.Time         `yaml:"mod"`   // Modification time
	Size           int64               `yaml:"size"`  // Filesize!
	ContentHash    map[string]string `yaml:"chash"` // Hash of the content
	PerceptualHash map[string]string `yaml:"phash"` // Hash of the image
}

// Generate : Generate new snapshot data from file, with all the trimmings. "err <-&Snapshot{}.Generate(name)"
func (sshot *Snapshot) Generate(filename string) chan error {
    done := make(chan error)
    go func(){
        var err error
        defer func(){
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
        phash, err := GeneratePerceptualHash("average", handle) // SHA256 hardcoded for now
        if err != nil {
            return
        }
        sshot.PerceptualHash = map[string]string{"average": phash}
    }()
    return done
}

// CheckFile : Check if a file matches corresponding snapshot
func CheckFile(filename string, sshot *Snapshot) (string, error) {
	// Get a handle
    handle, err := os.Open(filename)
    if err != nil {
        return "", err
    }
    defer handle.Close()
    info, err := handle.Stat()
    if err != nil {
        return "", err
    }

    // Some checking, escallating in complexity
    if filepath.Base(filename) != sshot.Name {
        return fmt.Sprintf("Name does not match '%s'", filename), nil
    }
    if info.Size() != sshot.Size {
        return fmt.Sprintf("Size does not match '%s'", filename), nil
    }
    if info.ModTime() == sshot.ModTime {
        // Roughly conclude a match!
        return "", nil
    }
    hash, err := GenerateContentHash("SHA256", handle) // SHA256 hardcoding for now
    if err != nil {
        return "", err
    }
    if hash != sshot.ContentHash["SHA256"] {
        return fmt.Sprintf("Content does not match '%s'", filename), nil
    }
    return "", nil
}

// TODO: manage file, listing snapshots

// todo validate files (is mod-time expected to be equal?) perhaps if size is same, and modtime is different, fall back to hash

// TODO: create lock function. perform lock on new files. validate existing files

// create "read only" function for linux/osx but also for windows
