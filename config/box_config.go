package config

//type ProjectRoot struct {
//	projectDir     string
//	configFilePath string
//}

type BoxConfig struct {
	HomeConfig    HomeConfig
	ProjectConfig ProjectConfig
}

func GetBoxConfig(homeConfig HomeConfig) BoxConfig {
	projectConfig := GetProjectConfig()
	return BoxConfig{
		ProjectConfig: projectConfig,
		HomeConfig:    homeConfig,
	}
}
