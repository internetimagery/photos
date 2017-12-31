package main

import (
	"fmt"
	"github.com/internetimagery/photos/cmd/initiate"
	"github.com/internetimagery/photos/cmd/rename"
	"github.com/internetimagery/photos/utility"
	"os"
	"strings"
)

func help() {
	fmt.Println("Shrink, Rename, Backup photos!")
	fmt.Println(">>>photos COMMAND ARGS")
	fmt.Println("(WIP) INIT   :: Set up the root of your photo project.")
	fmt.Println("(WIP) CONFIG :: Project settings")
	fmt.Println("(WIP) PROCESS:: Compress and rename photos.")
	fmt.Println("(WIP) BACKUP :: Copy files to another location.")
	fmt.Println("(WIP) DROP   :: Remove file from project, replacing with a pointer to original.")
	fmt.Println("(WIP) GET    :: Retrieve dropped file.")
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		help()
	} else {
		arg := strings.ToLower(os.Args[1])
		switch arg {
		case "init":
			initiate.Run(os.Args[2:])
		case "rename":
			rename.Run(os.Args[2:])
		default:
			guess := utility.ClosestMatch(arg, []string{"init", "rename"})
			fmt.Printf("Argument \"%s\" does not exist.\n", arg)
			fmt.Println("Did you mean:")
			fmt.Printf("\t%s", guess)
		}
	}
}
