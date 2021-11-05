package command

import (
	"fmt"
	"github.com/bseib/sideload/app"
	"github.com/bseib/sideload/config"
	"github.com/bseib/sideload/util"
	"os"
)

func Restore(sideloadConfig config.SideloadConfig, specificFiles []string, isForce bool) {
	if len(specificFiles) == 0 {
		comparisons := app.CompareProjectFiles(sideloadConfig)
		if len(comparisons.WouldRestore) == 0 {
			fmt.Printf("There are no tracked files that would restore.\n")
			os.Exit(1)
		} else {
			err := os.MkdirAll(sideloadConfig.ProjectConfig.ProjectDir, 0700)
			util.CheckFatal(err)
			for _, fileComparison := range comparisons.WouldRestore {
				_, err = app.CopyFile(fileComparison.HomeFile, fileComparison.ProjectFile)
				if err != nil {
					fmt.Printf("  error:  %s\n", fileComparison.RelativeFile)
					fmt.Printf("    %s\n", err)
				} else {
					fmt.Printf("  restored  %s\n", fileComparison.RelativeFile)
				}
			}
		}
	}
}
