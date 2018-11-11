package lock

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// TODO: generate phash
// TODO: make lock object to contain information
// TODO: impliment checking
// TODO: impliment serializing lock data
// TODO: make function to set files readonly, linux/osx/windows

// NewContentHash : Generate hash from content to compare contents
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
