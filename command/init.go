package command

import (
	"box/config"
	"flag"
)

func Init(flagset *flag.FlagSet) {
	config.InitProjectConfig()
}
