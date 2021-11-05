package command

import (
	"fmt"
	"github.com/bseib/sideload/app"
	"github.com/bseib/sideload/config"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"os"
)

func Diff(sideloadConfig config.SideloadConfig) {
	comparisons := app.CompareProjectFiles(sideloadConfig)
	if len(comparisons.AllComparisons) == 0 {
		fmt.Printf("There are no files being tracked by 'sideload'.\n")
		fmt.Printf("Edit '%v' to track some files.\n", config.CONFIG_FILENAME)
		return
	}
	if len(comparisons.WouldRestore) > 0 {
		fmt.Printf("Differences in files that 'sideload restore' would copy:\n")
		diffFileList(comparisons.WouldRestore, false)
		fmt.Println()
	}
	if len(comparisons.WouldStore) > 0 {
		fmt.Printf("Differences in files that 'sideload store' would copy:\n")
		diffFileList(comparisons.WouldStore, true)
		fmt.Println()
	}
}

func ternaryOfString(predicate bool, a string, b string) (result string) {
	if predicate {
		result = a
	} else {
		result = b
	}
	return
}

func diffFileList(files []app.FileComparison, swap bool) {
	for _, fileComparison := range files {
		fileA := ternaryOfString(swap, fileComparison.HomeFile, fileComparison.ProjectFile)
		fileB := ternaryOfString(swap, fileComparison.ProjectFile, fileComparison.HomeFile)
		fileAContent, err := readFileToString(fileA)
		if err != nil {
			fmt.Printf("  error:  %s\n", fileComparison.RelativeFile)
			fmt.Printf("    %s\n", err)
		} else {
			fileBContent, err := readFileToString(fileB)
			if err != nil {
				fmt.Printf("  error:  %s\n", fileComparison.RelativeFile)
				fmt.Printf("    %s\n", err)
			} else {
				edits := myers.ComputeEdits(span.URIFromPath(fileA), fileAContent, fileBContent)
				diff := fmt.Sprint(gotextdiff.ToUnified(fileA, fileB, fileAContent, edits))
				fmt.Print(diff)
			}
		}
	}
}

func readFileToString(file string) (string, error) {
	contents, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
