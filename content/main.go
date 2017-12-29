// Content hashing and comparison
package main

import (
  "fmt"
  "os"
  "io"
  "log"
  // "errors"
  "hash"
  "encoding/hex"
  "crypto/sha256"
  "crypto/sha1"
)



var PATH = "D:/Documents/go-workspace/src/github.com/internetimagery/photos/test.jpg"

// Return a hash representing the files content
func Hash(hash_type, path string) string {
  // Open file
  f, err := os.Open(path)
  if err != nil {
    log.Panic(err)
  }
  defer f.Close()

  // Choose hash
  var h hash.Hash

  switch hash_type {
    case "SHA1":
      h = sha1.New()
    case "SHA256":
      h = sha256.New()
    default:
      log.Fatalf("Hash type undefined %s", hash_type)
  }

  // Generate hash
  if _, err := io.Copy(h, f); err != nil {
    log.Panic(err)
  }

  // Return hash value
  return hex.EncodeToString(h.Sum(nil))
}

func main()  {
  fmt.Println(Hash("SHA1", PATH))
}
