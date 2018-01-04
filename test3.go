package main

import (
  "fmt"
  "flag"
  "os"
)

func main()  {
  countcmd := flag.NewFlagSet("count", flag.ExitOnError)

  countnum := countcmd.Int("num", 0, "counting number")

  countcmd.Parse(os.Args[1:])

  fmt.Println(os.Args)
  fmt.Println(*countnum)
  countcmd.PrintDefaults()
}
