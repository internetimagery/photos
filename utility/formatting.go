// File formatting utilities
package utility

import (
  "regexp"
)

// _num[tag tag].ext
const SUFFIX = "_(\\d+)(\\[(.+?)]|)\\.(\\w+)"

func GetRegex(prefix string) (*regexp.Regexp, error) {
  // Apply prefix and suffix. Return compiled regex
  return regexp.Compile(regexp.QuoteMeta(prefix) + SUFFIX)
}
