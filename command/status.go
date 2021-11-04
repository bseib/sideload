package command

import (
	"box/app"
	"box/config"
	"flag"
	"fmt"
)

func Status(flagset *flag.FlagSet, boxConfig config.BoxConfig) {
	comparisons := app.CompareProjectFiles(boxConfig)
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
		fmt.Println("Would restore:")
		printFileList(" -->", wouldRestore)
		fmt.Println()
	}
	if len(wouldStore) > 0 {
		fmt.Println("Would store:")
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
