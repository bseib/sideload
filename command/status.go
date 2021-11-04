package command

import (
	"flag"
	"fmt"
	"sideload/app"
	"sideload/config"
)

func Status(flagset *flag.FlagSet, sideloadConfig config.SideloadConfig) {
	comparisons := app.CompareProjectFiles(sideloadConfig)
	if len(comparisons) == 0 {
		fmt.Printf("There are no files being tracked by 'sideload'.\n")
		fmt.Printf("Edit '%v' to track some files.\n", config.CONFIG_FILENAME)
		return
	}

	wouldRestore := app.FilterOfFileComparison(comparisons, func(fc app.FileComparison) bool {
		return fc.Inclination == app.WILL_RESTORE
	})
	wouldStore := app.FilterOfFileComparison(comparisons, func(fc app.FileComparison) bool {
		return fc.Inclination == app.WILL_STORE
	})
	wouldNothing := app.FilterOfFileComparison(comparisons, func(fc app.FileComparison) bool {
		return fc.Inclination == app.NONE
	})
	bothMissing := app.FilterOfFileComparison(comparisons, func(fc app.FileComparison) bool {
		return fc.Inclination == app.BOTH_MISSING
	})

	if len(wouldRestore) > 0 {
		fmt.Printf("Would restore with 'sideload restore' from '%v':\n", sideloadConfig.HomeProjectDir)
		printFileList(" -->", wouldRestore)
		fmt.Println()
	}
	if len(wouldStore) > 0 {
		fmt.Printf("Would store with 'sideload store' into '%v':\n", sideloadConfig.HomeProjectDir)
		printFileList("<-- ", wouldStore)
		fmt.Println()
	}
	if len(wouldNothing) > 0 {
		fmt.Println("No changes:")
		printFileList(" == ", wouldNothing)
		fmt.Println()
	}
	if len(bothMissing) > 0 {
		fmt.Println("Missing:")
		printFileList(" !! ", bothMissing)
		fmt.Println()
	}
}

func printFileList(prefix string, files []app.FileComparison) {
	for _, file := range files {
		fmt.Printf("  %s  %s\n", prefix, file.Relativefile)
	}
}
