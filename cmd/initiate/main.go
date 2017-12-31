package initiate

import (
  "log"
  "fmt"
  "os"
  "path/filepath"
  "github.com/internetimagery/photos/config"
)

func help()  {
  fmt.Println("Initiate new repo")
  fmt.Println(">>>photos INIT <name>")
}

func Run(args []string)  {
  if len(args) != 0 && (args[0] == "-h" || args[1] == "--help"){
    help()
  } else {
    // Check if config file already exists.
    cwd, err := os.Getwd()
    if err != nil {
      log.Panic(err)
    }
    if config.GetConfig(cwd) == "" {
      log.Println("Initializing new Repo.")
      conf := config.NewConfig()
      path := filepath.Join(cwd, config.CONFIGNAME)
      fmt.Println(conf)
      fmt.Println(path)

    } else {
      log.Panic("Already Initialized.")
    }
  }
}
