package command

import (
	"fmt"
	"github.com/bseib/sideload/app"
	"github.com/bseib/sideload/config"
	"github.com/bseib/sideload/util"
	"os"
)

func Store(sideloadConfig config.SideloadConfig, specificFiles []string, isForce bool) {
	if len(specificFiles) == 0 {
		comparisons := app.CompareProjectFiles(sideloadConfig)
		if len(comparisons.WouldStore) == 0 {
			fmt.Printf("There are no tracked files that would store.\n")
			os.Exit(1)
		} else {
			err := os.MkdirAll(sideloadConfig.HomeProjectDir, 0700)
			util.CheckFatal(err)
			for _, fileComparison := range comparisons.WouldStore {
				_, err = app.CopyFile(fileComparison.ProjectFile, fileComparison.HomeFile)
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
