// Content hashing and comparison
// package main
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

// Return number (string) representing filesize
func Size(file *os.File) string {
  // Get file stats
  stat, err := file.Stat()
  if err != nil {
    log.Panic(err)
  }
  return strconv.FormatInt(stat.Size(), 10)
}

// Return a hash representing the files content
func Hash(hash_type string, file *os.File) string {
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
  if _, err := io.Copy(h, file); err != nil {
    log.Panic(err)
  }
  return hex.EncodeToString(h.Sum(nil))
}

// func main()  {
//   p := "D:/Documents/go-workspace/src/github.com/internetimagery/photos/test.jpg"
//   // Open file
//   f, err := os.Open(p)
//   if err != nil {
//     log.Panic(err)
//   }
//   defer f.Close()
//
//   fmt.Println(Hash("SHA1", f))
// }
