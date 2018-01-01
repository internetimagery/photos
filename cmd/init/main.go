package cmdinit

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
  numArgs := len(args)
  cwd, err := os.Getwd()
  if err != nil {
    log.Panic(err)
  }
  if numArgs != 0 && (args[0] == "-h" || args[0] == "--help"){
    help()
  } else {
    // Check if config file already exists.
    if config.GetConfig(cwd) == "" {
      log.Println("Initializing new Repo.")
      conf := config.NewConfig()
      path := filepath.Join(cwd, config.CONFIGNAME)
      if numArgs == 1 {
        conf.Name = args[0]
      }
      config.SaveConfig(conf, path)
    } else {
      log.Println("Already Initialized.")
    }
  }
}
