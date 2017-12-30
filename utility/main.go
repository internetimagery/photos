// General utility
package utility

import (
  "os"
  "net"
  "fmt"
  "log"
  "bytes"
  "time"
  "crypto/sha1"
  // "encoding/base64"
  "github.com/satori/go.uuid"
  // "crypto/rand"
)

// Get current working directory
func CWD() string {
  cwd, err := os.Getwd()
  if err != nil {
    log.Fatal(err)
  }
  return cwd
}

// Get mac address
func GetMAC() string {
  interfaces, err := net.Interfaces()
  if err != nil {
    log.Fatal(err)
  }
  for _, i := range interfaces {
    if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
      addr := i.HardwareAddr.String()
      return addr
    }
  }
  return ""
}

// RandomID from mac address and time
func GenerateID() string {
  return uuid.NewV4()
  // mac := GetMAC()
  // time := time.Now().String()
  // hash := sha1.New().Sum([]byte(time+mac))
  // return base64.StdEncoding.EncodeToString(hash)
}
