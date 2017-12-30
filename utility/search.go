// walk file path
package utility

import (
  "io/ioutil"
  "path/filepath"
)

func SearchUp(search, root string) string {
  for true {
    files, _ := ioutil.ReadDir(root)
    for i := 0; i < len(files); i++ {
      name := files[i].Name()
      if search == name {
        return filepath.Join(root, name)
      }
    }
    new_root := filepath.Dir(root)
    if new_root == root {
      break
    }
    root = new_root
  }
  return ""
}
