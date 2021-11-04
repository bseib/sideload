package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type HomeConfig struct {
	HomeDir    string
	StorageDir string
}

func getSideloadHomeDirEnv() string {
	return os.Getenv("SIDELOADHOME")
}

func getSideloadHomeDirPath() string {
	var sideloadHomeDirEnv = getSideloadHomeDirEnv()
	if len(sideloadHomeDirEnv) == 0 {
		homedir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Cannot determine your home directory.")
			log.Fatal(err)
		}
		return filepath.Join(filepath.Clean(homedir), ".sideload")
	} else {
		return filepath.Clean(sideloadHomeDirEnv)
	}
}

func GetHomeConfig() HomeConfig {
	sideloadHomeDirPath := getSideloadHomeDirPath()
	sideloadStorageDir := filepath.Join(filepath.Clean(sideloadHomeDirPath), "storage")
	if _, err := os.Stat(sideloadHomeDirPath); os.IsNotExist(err) {
		var sideloadHomeDirEnv = getSideloadHomeDirEnv()
		if len(sideloadHomeDirEnv) == 0 {
			fmt.Printf("SIDELOADHOME environment variable has not been set. Using default: '%v'\n", sideloadHomeDirPath)
		}
		fmt.Printf("Directory '%v' does not exist. Create it now? [Y/n] ", sideloadHomeDirPath)
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		text = strings.TrimRight(text, "\r\n")
		if err == nil && (len(text) == 0 || strings.HasPrefix(text, "y") || strings.HasPrefix(text, "Y")) {
			fmt.Printf("Creating directory '%v'.\n", sideloadHomeDirPath)
			err = os.Mkdir(sideloadHomeDirPath, 0700)
			if err != nil {
				log.Fatal(err)
			}
			err = os.Mkdir(sideloadStorageDir, 0700)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("A SIDELOADHOME directory is required and must exist to do anything.")
			os.Exit(1)
		}
	}
	return HomeConfig{
		HomeDir:    sideloadHomeDirPath,
		StorageDir: sideloadStorageDir,
	}
}
