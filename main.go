package main

import (
	"fmt"
	"github.com/internetimagery/photos/rename"
	"os"
	"strings"
)

func help() {
	fmt.Println("Shrink, Rename, Backup photos!")
	fmt.Println("photos COMMAND ARGS")
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
		switch strings.ToLower(os.Args[1]) {
		case "init":
			rename.Run(os.Args[2:])
		case "rename":
			rename.Run(os.Args[2:])
		default:
			help()
		}
	}
}
