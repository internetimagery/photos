package main

import (
	"fmt"
	"github.com/internetimagery/photos/cmd/init"
	"github.com/internetimagery/photos/cmd/config"
	"github.com/internetimagery/photos/utility"
	"github.com/internetimagery/photos/config"
	"os"
	"log"
	"strings"
)

type run func([]string, *config.Config)

var ARGS = map[string]run{
	"init": cmdinit.Run,
	"config": cmdconfig.Run,
}

func help() {
	fmt.Println("Shrink, Rename, Backup photos!")
	fmt.Println(">>>photos COMMAND ARGS")
	fmt.Println("(WIP) INIT   :: Set up the root of your photo project.")
	fmt.Println("(WIP) CONFIG :: Project settings")
	fmt.Println("(WIP) ADD    :: Compress and rename photos.")
	fmt.Println("(WIP) BACKUP :: Copy files to another location.")
	fmt.Println("(WIP) DROP   :: Remove file from project, replacing with a pointer to original.")
	fmt.Println("(WIP) GET    :: Retrieve dropped file.")
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		help()
	} else {
		arg := strings.ToLower(os.Args[1])
		cwd, err := os.Getwd()
		if err != nil {
			log.Panic(err)
		}
		conf := config.GetConfig(cwd)
		if arg == "init" {
			cmdinit.Run(conf, conf)
		}
		if val, ok := ARGS[arg]; ok {
			if conf.GetRoot() != "" || arg == "init" {
				val(os.Args[2:], conf)
				return
			}
			log.Panic("Not inside repository")
		} else {
			options := make([]string, len(ARGS))
			i := 0
			for k := range ARGS {
				options[i] = k
				i++
			}
			guess := utility.ClosestMatch(arg, options)
			fmt.Printf("Argument \"%s\" does not exist.\n", arg)
			fmt.Println("Did you mean:")
			fmt.Printf("\t%s", guess)
		}
	}
}
