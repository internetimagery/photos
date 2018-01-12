package cmdinit

import (
  "log"
  "fmt"
  // "os"
  // "path/filepath"
  "github.com/internetimagery/photos/config"
)

func help()  {
  fmt.Println("Initiate new repo")
  fmt.Println(">>>photos INIT <name>")
}

type Command struct{}

func (_ *Command) Desc() string {
  return "HI THERE"
}

func (_ *Command) Run(args []string, conf *config.Config) {
  log.Println("Running init")
  log.Println(args)
}

// func Run(args []string, conf *config.Config)  {
//   numArgs := len(args)
//   if numArgs != 0 && (args[0] == "-h" || args[0] == "--help"){
//     help()
//   } else {
//     // Check if config file already exists.
//     if conf.GetRoot() == "" {
//       log.Println("Initializing new Repo.")
//       cwd, _ := os.Getwd()
//       path := filepath.Join(cwd, config.CONFIGNAME)
//       if numArgs == 1 {
//         conf.Name = args[0]
//       }
//       config.SaveConfig(conf, path)
//     } else {
//       log.Println("Already Initialized.")
//     }
//   }
// }
