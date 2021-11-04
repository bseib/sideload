package config

import "path/filepath"

type SideloadConfig struct {
	HomeConfig     HomeConfig
	ProjectConfig  ProjectConfig
	HomeProjectDir string
}

func GetSideloadConfig(homeConfig HomeConfig) SideloadConfig {
	projectConfig := GetProjectConfig()
	return SideloadConfig{
		ProjectConfig:  projectConfig,
		HomeConfig:     homeConfig,
		HomeProjectDir: filepath.Join(homeConfig.StorageDir, projectConfig.Project.Name),
	}
}
