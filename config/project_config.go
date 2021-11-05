package config

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/bseib/sideload/util"
	"os"
	"path/filepath"
	"strings"
)

const CONFIG_FILENAME = ".sideload-config"

type Project struct {
	Name string
}

type Files struct {
	Track []string
}

type ProjectConfigToml struct {
	Project Project `toml:"project"`
	Files   Files   `toml:"files"`
}

type ProjectConfig struct {
	ProjectDir string
	Project    Project
	Files      Files
}

func getCurrentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		util.Fatal(err)
	}
	return currentDir
}

func findConfigFilePath(dir string) (string, error) {
	lookForFile := filepath.Join(dir, CONFIG_FILENAME)
	if _, err := os.Stat(lookForFile); !os.IsNotExist(err) {
		// found a sideload config file
		return lookForFile, nil
	}
	parentDir := filepath.Clean(filepath.Join(dir, ".."))
	if parentDir == dir {
		return "", errors.New(fmt.Sprintf("Could not find a '%s' file on a walk from this directory up to the root directory.", CONFIG_FILENAME))
	} else {
		return findConfigFilePath(parentDir)
	}
}

func InitProjectConfig() ProjectConfig {
	path, err := findConfigFilePath(getCurrentDir())
	if err == nil {
		fmt.Printf("This directory is already covered by sideload with config file '%s'.\n", path)
		os.Exit(1)
	}
	fmt.Printf("Config file '%v' does not exist. Create it now? [Y/n] ", CONFIG_FILENAME)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	text = strings.TrimRight(text, "\r\n")
	if err == nil && (len(text) == 0 || strings.HasPrefix(text, "y") || strings.HasPrefix(text, "Y")) {
		fmt.Printf("Creating config file '%v'.\n", CONFIG_FILENAME)
		WriteDefaultConfig()
	} else {
		fmt.Printf("Okay, no '%v' was file created.\n", CONFIG_FILENAME)
		os.Exit(1)
	}
	return GetProjectConfig()
}

func GetProjectConfig() ProjectConfig {
	path, err := findConfigFilePath(getCurrentDir())
	if err != nil {
		fmt.Printf("This directory has not been initialized to use sideload. Try 'sideload init' to create a '%s' file.\n", CONFIG_FILENAME)
		os.Exit(1)
	}
	data, err := os.ReadFile(path)
	tomlConfig := ProjectConfigToml{}
	_, err = toml.Decode(string(data), &tomlConfig)
	if err != nil {
		fmt.Printf("Error parsing file '%v':\n", path)
		fmt.Println(err)
		os.Exit(1)
	}

	config := ProjectConfig{
		ProjectDir: filepath.Dir(path),
		Project:    tomlConfig.Project,
		Files:      tomlConfig.Files,
	}
	if len(config.Project.Name) == 0 {
		util.Fatal(fmt.Sprintf("Project name cannot be empty in '%s'.", path))
	}
	if config.Files.Track == nil {
		config.Files.Track = []string{}
	}
	return config
}

func WriteDefaultConfig() {
	currentDir := filepath.Base(getCurrentDir())
	writeConfig(currentDir, []string{})
}

func writeConfig(projectName string, files []string) {
	fileList := strings.Join(util.MapOfString(files, func(s string) string {
		return fmt.Sprintf("  '%v'", s)
	}), ",")
	defaultConfigDoc := []byte(fmt.Sprintf(`[project]
##
## This is a human readable name representing this set of tracked files.
##
name = '%s'

[files]
##
## List the files you want to track when you sideload files.
## File names are relative to this directory. Example:
##
## track = [
##   'filename1',
##   'dir/filename2',
## ]
##
track = [
%s
]

`, projectName /*projectId,*/, fileList))
	err := os.WriteFile(CONFIG_FILENAME, defaultConfigDoc, 0644)
	if err != nil {
		util.Fatal(err)
	}
}
