package command

import (
	"fmt"
	"github.com/bseib/sideload/app"
	"github.com/bseib/sideload/config"
)

func Status(sideloadConfig config.SideloadConfig) {
	comparisons := app.CompareProjectFiles(sideloadConfig)
	if len(comparisons.AllComparisons) == 0 {
		fmt.Printf("There are no files being tracked by 'sideload'.\n")
		fmt.Printf("Edit '%v' to track some files.\n", config.CONFIG_FILENAME)
		return
	}
	if len(comparisons.WouldRestore) > 0 {
		fmt.Printf("Would restore with 'sideload restore' from '%v':\n", sideloadConfig.HomeProjectDir)
		printFileList(" -->", comparisons.WouldRestore)
		fmt.Println()
	}
	if len(comparisons.WouldStore) > 0 {
		fmt.Printf("Would store with 'sideload store' into '%v':\n", sideloadConfig.HomeProjectDir)
		printFileList("<-- ", comparisons.WouldStore)
		fmt.Println()
	}
	if len(comparisons.WouldNothing) > 0 {
		fmt.Println("No changes:")
		printFileList(" == ", comparisons.WouldNothing)
		fmt.Println()
	}
	if len(comparisons.BothMissing) > 0 {
		fmt.Println("Missing:")
		printFileList(" !! ", comparisons.BothMissing)
		fmt.Println()
	}
}

func printFileList(prefix string, files []app.FileComparison) {
	for _, file := range files {
		fmt.Printf("  %s  %s\n", prefix, file.RelativeFile)
	}
}
