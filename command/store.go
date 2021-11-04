package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sideload/app"
	"sideload/config"
	"sideload/util"
)

func copyFile(src string, dest string) (bytesWritten int64, err error) {
	err = os.MkdirAll(filepath.Dir(dest), 0700)
	if err != nil {
		return 0, err
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()
	destFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()
	bytesWritten, err = io.Copy(destFile, srcFile)
	if err != nil {
		return bytesWritten, err
	}
	return bytesWritten, nil
}

func Store(sideloadConfig config.SideloadConfig, specificFiles []string, isForce bool) {
	if len(specificFiles) == 0 {
		comparisons := app.CompareProjectFiles(sideloadConfig)
		if len(comparisons.WouldStore) == 0 {
			fmt.Printf("There are no tracked files to store.\n")
			os.Exit(1)
		} else {
			err := os.MkdirAll(sideloadConfig.HomeProjectDir, 0700)
			util.CheckFatal(err)
			for _, fileComparison := range comparisons.WouldStore {
				_, err = copyFile(fileComparison.ProjectFile, fileComparison.HomeFile)
				if err != nil {
					fmt.Printf("  error:  %s\n", fileComparison.RelativeFile)
					fmt.Printf("    %s\n", err)
				} else {
					fmt.Printf("  stored  %s\n", fileComparison.RelativeFile)
				}
			}
		}
	}
}
