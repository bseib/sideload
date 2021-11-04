package main

import (
	"flag"
	"fmt"
	"os"
	"sideload/command"
	"sideload/config"
)

func main() {
	var homeConfig = config.GetHomeConfig()

	// Subcommand flags
	initFlagSet := flag.NewFlagSet("init", flag.ExitOnError)
	storeFlagSet := flag.NewFlagSet("store", flag.ExitOnError)
	var storeForce = storeFlagSet.Bool("f", false, "Force specified files to be stored")

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
			command.Status(config.GetSideloadConfig(homeConfig))
		case "store":
			storeFlagSet.Parse(os.Args[2:])
			command.Store(config.GetSideloadConfig(homeConfig), storeFlagSet.Args(), *storeForce)
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
		"   init      init a directory to manage its sideloaded (tracked) files\n" +
		"   status    show which tracked files would be copied to/from storage\n" +
		"   store     store tracked files to storage dir, if they are newer\n" +
		"   restore   restore tracked files to local project dir from storage, if they are newer\n")
	flag.PrintDefaults()
	os.Exit(1)
}
