// Testing formatting
package utility

import (
  "testing"
)


func TestGetRegex(t *testing.T) {
  prefix := "one two three"
  exp, err := GetRegex(prefix)
  if err != nil {
    t.FailNow()
  }
  if !exp.MatchString(prefix + "_001[one two].ext1") {
    t.Fail()
  }
  if !exp.MatchString(prefix + "_021.mp4") {
    t.Fail()
  }
  if exp.MatchString(prefix + ".mp4") {
    t.Fail()
  }
  if exp.MatchString("one two four.mp4") {
    t.Fail()
  }
}

func TestParseFiles(t *testing.T) {
  prefix := "testing 123"
  files := []string {
    prefix + "stuff",
    prefix + "_002[tag1 tag2].jpg",
    "something else",
    prefix + "_034" }
  media, err := ParseFiles(prefix, files)
  if err == nil {
    t.FailNow()
  }
}
