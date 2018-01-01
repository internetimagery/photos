package main

import (
	"fmt"
	"github.com/internetimagery/photos/cmd/init"
	"github.com/internetimagery/photos/cmd/config"
	"github.com/internetimagery/photos/utility"
	"os"
	"strings"
)

type run func([]string)

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
		if val, ok := ARGS[arg]; ok {
			val(os.Args[2:])
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
