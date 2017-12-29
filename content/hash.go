// Content hashing and comparison
package content

import (
  "fmt"
  "os"
  "io"
  "log"
  "hash"
  "encoding/hex"
  "crypto/sha256"
  "crypto/sha1"
  "strconv"
)


// Return a hash representing the files content
func Hash(hash_type, path string) string {
  // Open file
  f, err := os.Open(path)
  if err != nil {
    log.Panic(err)
  }
  defer f.Close()

  // Get file stats
  stat, err := f.Stat()
  if err != nil {
    log.Panic(err)
  }
  size := strconv.FormatInt(stat.Size(), 10)

  // Choose hash
  var h hash.Hash
  switch hash_type {
    case "SHA1":
      h = sha1.New()
    case "SHA256":
      h = sha256.New()
    default:
      log.Fatalf("Hash type undefined: %s", hash_type)
  }

  // Generate hash
  if _, err := io.Copy(h, f); err != nil {
    log.Panic(err)
  }
  fingerprint := hex.EncodeToString(h.Sum(nil))

  // Return hash value
  return size + "-" + fingerprint
}
