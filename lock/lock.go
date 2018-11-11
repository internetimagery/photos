package lock

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io"

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
		return base64.StdEncoding.EncodeToString(hasher.Sum([]byte{})), nil
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
