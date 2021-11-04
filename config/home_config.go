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
	homeDirPath string
}

func getBoxHomeDirEnv() string {
	return os.Getenv("BOXHOME")
}

func getBoxHomeDirPath() string {
	var boxHomeDirEnv = getBoxHomeDirEnv()
	if len(boxHomeDirEnv) == 0 {
		homedir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Cannot determine your home directory.")
			log.Fatal(err)
		}
		return filepath.Join(filepath.Clean(homedir), ".box")
	} else {
		return filepath.Clean(boxHomeDirEnv)
	}
}

func GetHomeConfig() HomeConfig {
	boxHomeDirPath := getBoxHomeDirPath()
	if _, err := os.Stat(boxHomeDirPath); os.IsNotExist(err) {
		var boxHomeDirEnv = getBoxHomeDirEnv()
		if len(boxHomeDirEnv) == 0 {
			fmt.Printf("BOXHOME environment variable has not been set. Using default: '%v'\n", boxHomeDirPath)
		}
		fmt.Printf("Directory '%v' does not exist. Create it now? [Y/n] ", boxHomeDirPath)
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		text = strings.TrimRight(text, "\r\n")
		if err == nil && (len(text) == 0 || strings.HasPrefix(text, "y") || strings.HasPrefix(text, "Y")) {
			fmt.Printf("Creating directory '%v'.\n", boxHomeDirPath)
			err = os.Mkdir(boxHomeDirPath, 0700)
			if err != nil {
				log.Fatal(err)
			}
			boxStorageDir := filepath.Join(filepath.Clean(boxHomeDirPath), "storage")
			err = os.Mkdir(boxStorageDir, 0700)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("A BOXHOME directory is required and must exist to do anything.")
			os.Exit(1)
		}
	}
	return HomeConfig{
		homeDirPath: boxHomeDirPath,
	}
}
