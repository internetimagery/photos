// Generate / compare file hashes
package utility

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
)

// Generate Hash from a file
func GetHashFromFile(filename, hashtype string) (string, error) {
	var hasher hash.Hash
	Fhandle, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer Fhandle.Close()

	// Pick hasher
	switch hashtype {
	case "md5":
		hasher = md5.New()
	default:
		return "", errors.New(fmt.Sprintf("Unrecognized hashtype %s", hashtype))
	}

	// Digest hash
	if _, err := io.Copy(hasher, Fhandle); err != nil {
		return "", err
	}

	result := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	// result := hex.EncodeToString(hasher.Sum(nil))
	return fmt.Sprintf("%s-%s", hashtype, result), nil
}
