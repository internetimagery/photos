// File utilities.

package utility

import (
  "gopkg.in/h2non/filetype.v1"
  "io/ioutil"
  "os"
)

const UNKNOWN = 0
const IMAGE = 1
const VIDEO = 2

func getHeader(path string) []byte {
  file, err := os.Open(path)
  if err != nil {
    Panic(err)
  }
  defer os.Close
  header := make([]byte, 261)
  file.Read(header)
  return file
}

// Get file type
func GetFileType(path string) int {
  header := getHeader(path)
  switch true {
  case filetype.IsImage(header):
    return IMAGE
  case filetype.IsVideo(header):
    return VIDEO
  default:
    return UNKNOWN
  }
}
