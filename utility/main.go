// General utility
package utility

import (
  "os"
  // "net"
  // "fmt"
  "log"
  // "bytes"
  // "time"
  // "crypto/sha1"
  "github.com/satori/go.uuid"
  "github.com/schollz/closestmatch"
)

func MaxInt(a, b int) int {
  if a > b {
    return a
  }
  return b
}

// Get current working directory
func CWD() string {
  cwd, err := os.Getwd()
  if err != nil {
    log.Panic(err)
  }
  return cwd
}

// Get closest match and return it
func ClosestMatch(word string, list []string) string {
  cm := closestmatch.New(list, []int{2})
  return cm.Closest(word)
}

// // Get mac address
// func GetMAC() string {
//   interfaces, err := net.Interfaces()
//   if err != nil {
//     log.Fatal(err)
//   }
//   for _, i := range interfaces {
//     if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
//       addr := i.HardwareAddr.String()
//       return addr
//     }
//   }
//   return ""
// }

// RandomID from mac address and time
func GenerateID() string {
  return uuid.NewV4().String()
}
