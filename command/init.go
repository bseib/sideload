package command

import (
	"github.com/bseib/sideload/config"
)

func Init(homeConfig config.HomeConfig) {
	config.InitProjectConfig()
}
