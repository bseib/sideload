package config

import (
	"box/util"
	"bufio"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const CONFIG_FILENAME = ".box-config"

type Project struct {
	Name string
	Uuid string
}

type Files struct {
	Track []string
}

type ProjectConfig struct {
	Project Project `toml:"project"`
	Files   Files   `toml:"files"`
}

func getCurrentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return currentDir
}

func findConfigFilePath(dir string) (string, error) {
	lookForFile := filepath.Join(dir, CONFIG_FILENAME)
	if _, err := os.Stat(lookForFile); !os.IsNotExist(err) {
		// found a box config file
		return lookForFile, nil
	}
	parentDir := filepath.Clean(filepath.Join(dir, ".."))
	if parentDir == dir {
		return "", errors.New(fmt.Sprintf("Could not find a '%s' file on a walk from this directory up to the root directory.", CONFIG_FILENAME))
	} else {
		return findConfigFilePath(parentDir)
	}
}

func InitProjectConfig() {
	path, err := findConfigFilePath(getCurrentDir())
	if err == nil {
		fmt.Printf("This directory is already covered by box with config file '%s'.\n", path)
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
}

func GetProjectConfig() ProjectConfig {
	path, err := findConfigFilePath(getCurrentDir())
	if err != nil {
		fmt.Printf("This directory has not been initialized to use box. Try 'box init' to create a '%s' file.\n", CONFIG_FILENAME)
		os.Exit(1)
	}
	//fmt.Printf("filepath=%v\n", path)
	data, err := os.ReadFile(path)
	//fmt.Println(string(data))
	config := ProjectConfig{}
	_, err = toml.Decode(string(data), &config)
	util.CheckFatal(err)

	if len(config.Project.Name) == 0 {
		util.Fatal(fmt.Sprintf("Project name cannot be empty in '%s'.", path))
	}
	if len(config.Project.Uuid) == 0 {
		util.Fatal(fmt.Sprintf("Project uuid cannot be empty in '%s'.", path))
	}
	if config.Files.Track == nil {
		config.Files.Track = []string{}
	}
	return config
}

func WriteDefaultConfig() {
	currentDir := filepath.Base(getCurrentDir())
	projectId := uuid.New().String()
	writeConfig(currentDir, projectId, []string{})
}

func writeConfig(projectName string, projectId string, files []string) {
	fileList := strings.Join(util.MapOfString(files, func(s string) string {
		return fmt.Sprintf("  '%v'", s)
	}), ",")
	defaultConfigDoc := []byte(fmt.Sprintf(`[project]
##
## This is a human readable name (for this uuid) representing this set of tracked files.
##
name = '%s'

##
## Tracked files for this project ultimately get stored in $BOXHOME, under this uuid.
##
uuid = '%s'

[files]
##
## List the files you want to track when you box/unbox files.
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

`, projectName, projectId, fileList))
	err := os.WriteFile(CONFIG_FILENAME, defaultConfigDoc, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
