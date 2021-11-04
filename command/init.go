package command

import (
	"sideload/config"
)

func Init(homeConfig config.HomeConfig) {
	config.InitProjectConfig()
}
