package command

import (
	"sideload/config"
	"flag"
)

func Init(flagset *flag.FlagSet, homeConfig config.HomeConfig) {
	config.InitProjectConfig()
}
