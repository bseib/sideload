package config

//type ProjectRoot struct {
//	projectDir     string
//	configFilePath string
//}

type BoxConfig struct {
	homeConfig    HomeConfig
	projectConfig ProjectConfig
}

func GetBoxConfig(homeConfig HomeConfig) BoxConfig {
	projectConfig := GetProjectConfig()
	return BoxConfig{
		projectConfig: projectConfig,
		homeConfig:    homeConfig,
	}
}
