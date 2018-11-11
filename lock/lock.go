package lock

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io"
	"time"

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
	case "MD5":
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
	case "Average":
		hash, err := goimagehash.AverageHash(img)
		if err != nil {
			return "", err
		}
		return hash.ToString(), nil
	}
	return "", fmt.Errorf("Unknown hash format '%s'", hashType)
}

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
	Size           int               `yaml:"size"`  // Filesize!
	ContentHash    map[string]string `yaml:"chash"` // Hash of the content
	PerceptualHash map[string]string `yaml:"phash"` // Hash of the image
}

// TODO: new snapshot, create all data

// TODO: manage file, listing snapshots

// todo validate files (is mod-time expected to be equal?) perhaps if size is same, and modtime is different, fall back to hash

// TODO: create lock function. perform lock on new files. validate existing files

// create "read only" function for linux/osx but also for windows
