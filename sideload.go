package main

import (
	"sideload/command"
	"sideload/config"
	"flag"
	"fmt"
	"os"
)

func main() {
	var homeConfig = config.GetHomeConfig()

	// Subcommands
	initFlagSet := flag.NewFlagSet("init", flag.ExitOnError)
	statusFlagSet := flag.NewFlagSet("status", flag.ExitOnError)

	// init flags
	//	var initYo = initFlagSet.String("yo", "", "yo yo yo")
	//	var statusYo = statusFlagSet.String("yo", "", "yo yo yo")

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		dieUsage()
	} else {
		switch os.Args[1] {
		case "init":
			initFlagSet.Parse(os.Args[2:])
			command.Init(initFlagSet, homeConfig)
		case "restore":
			command.Restore()
		case "status":
			statusFlagSet.Parse(os.Args[2:])
			command.Status(statusFlagSet, config.GetSideloadConfig(homeConfig))
		case "store":
			command.Store()
		default:
			dieUsage()
		}
	}

	//var nFlag = flag.Int("n", 1234, "help message for flag n")
	//flag.Parse()
	//fmt.Printf("Hello, nFlag=%d\n", *nFlag)

}

func dieUsage() {
	fmt.Println("usage: sideload <command> [<args>]\n" +
		"\n" +
		"Common commands:\n" +
		"   init      init a directory to manage its sideloaded files\n" +
		"   status    show which sideloaded files need copied or have changed\n")
	flag.PrintDefaults()
	os.Exit(1)
}
