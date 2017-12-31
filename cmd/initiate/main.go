package initiate

import (
  "fmt"
)

func help()  {
  fmt.Println("Initiate new repo")
  fmt.Println("photos INIT")
}

func Run(args []string)  {
  if len(args) != 0 {
    help()
  } else {
    fmt.Println("INIT")
    fmt.Println(args)
  }
}
