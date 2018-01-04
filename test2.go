package main

import (
  "fmt"
  "github.com/internetimagery/photos/config"
)

func main()  {
  conf := config.NewConfig()
  conf.Name = "testingggg"
  p := "D:/Documents/go-workspace/src/github.com/internetimagery/photos/conf.json"
  config.SaveConfig(conf, p)
  fmt.Println(config.LoadConfig(p))
}
