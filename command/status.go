package command

import (
	"box/config"
	"flag"
	"fmt"
)

func Status(flagset *flag.FlagSet, boxConfig config.BoxConfig) {
	config.GetProjectConfig()
	fmt.Println("TODO")
}