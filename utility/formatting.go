// File formatting utilities
package utility

import (
  "regexp"
)

// _num[tag tag].ext
const SUFFIX = "_(\\d+)(?:\\[(.+?)])?\\.(\\w+)"

func GetRegex(prefix string) (*regexp.Regexp, error) {
  // Apply prefix and suffix. Return compiled regex
  return regexp.Compile(regexp.QuoteMeta(prefix) + SUFFIX)
}

type Media struct {
  Prefix, Ext string
  Index int
  Tags []string
}

func ParseFiles(parent string, files []string) ([]Media, error){
  reg, err := GetRegex(parent)
  result := []Media{}
  for i := 0; i < len(files); i++ {
    result.append(Media{})
  }
  return result, err
}
