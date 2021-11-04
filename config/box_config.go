package config

//type ProjectRoot struct {
//	projectDir     string
//	configFilePath string
//}

type SideloadConfig struct {
	HomeConfig    HomeConfig
	ProjectConfig ProjectConfig
}

func GetSideloadConfig(homeConfig HomeConfig) SideloadConfig {
	projectConfig := GetProjectConfig()
	return SideloadConfig{
		ProjectConfig: projectConfig,
		HomeConfig:    homeConfig,
	}
}
