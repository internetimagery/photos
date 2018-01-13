package main

import (
	"fmt"
	// "github.com/internetimagery/photos/cmd/init"
	// "github.com/internetimagery/photos/cmd/config"
	"github.com/internetimagery/photos/utility"
	"github.com/internetimagery/photos/state"
	"github.com/internetimagery/photos/commands"
	"os"
	// "log"
	"strings"
)

// Simple command
type Command interface {
	Desc() string
	Run([]string, *state.State) int
}

// Brief help message
func help(mod map[string]Command)  {
	fmt.Println("Usage:\n")
	max := 0
	for name, _ := range mod {
		max = utility.MaxInt(max, len(name))
	}
	for name, com := range mod {
		fmt.Println(strings.Repeat(" ", max-len(name)) + name, ":", com.Desc())
	}
}

// Lets go!
func main()  {
	// No arguments? Show help message
	modules := map[string]Command{
		"init": commands.CMD_Init{},
		"config": commands.CMD_Config{},
		"add": commands.CMD_Add{},
		"drop": commands.CMD_Drop{},
		"get": commands.CMD_Get{},
		"backup": commands.CMD_Backup{},
	}
	if len(os.Args) < 2 {
		help(modules)
	} else {
		fnc, ok := modules[strings.ToLower(os.Args[1])]
		if !ok {
			help(modules)
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				log.Panic(err)
			}
			state := state.State.New(cwd)
			os.Exit(fnc.Run(os.Args[2:], state))
		}
	}
}

//
// func main() {
// 	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
// 		help()
// 	} else {
// 		arg := strings.ToLower(os.Args[1])
// 		cwd, err := os.Getwd()
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		conf := config.GetConfig(cwd)
// 		if arg == "init" {
// 			cmdinit.Run(conf, conf)
// 		}
// 		if val, ok := ARGS[arg]; ok {
// 			if conf.GetRoot() != "" || arg == "init" {
// 				val(os.Args[2:], conf)
// 				return
// 			}
// 			log.Panic("Not inside repository")
// 		} else {
// 			options := make([]string, len(ARGS))
// 			i := 0
// 			for k := range ARGS {
// 				options[i] = k
// 				i++
// 			}
// 			guess := utility.ClosestMatch(arg, options)
// 			fmt.Printf("Argument \"%s\" does not exist.\n", arg)
// 			fmt.Println("Did you mean:")
// 			fmt.Printf("\t%s", guess)
// 		}
// 	}
// }
