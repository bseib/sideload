package app

import (
	"crypto/md5"
	"fmt"
	"github.com/bseib/sideload/config"
	"github.com/bseib/sideload/util"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type InclinationEnum int

const (
	NONE         InclinationEnum = iota
	WILL_STORE   InclinationEnum = iota
	WILL_RESTORE InclinationEnum = iota
	BOTH_MISSING InclinationEnum = iota
)

type FileComparison struct {
	RelativeFile string
	HomeFile     string
	ProjectFile  string
	Inclination  InclinationEnum
}

type ProjectFilesComparison struct {
	AllComparisons []FileComparison
	WouldRestore   []FileComparison
	WouldStore     []FileComparison
	WouldNothing   []FileComparison
	BothMissing    []FileComparison
}

type FileSituation struct {
	exists  bool
	modTime time.Time
	md5Hash string
}

func md5Hash(file string) string {
	f, err := os.Open(file)
	if err != nil {
		util.Fatal(err)
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		util.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func getFileSituation(file string) FileSituation {
	stat, err := os.Stat(file)
	if os.IsNotExist(err) {
		return FileSituation{
			exists:  false,
			modTime: time.Time{},
			md5Hash: "",
		}
	} else {
		return FileSituation{
			exists:  true,
			modTime: stat.ModTime(),
			md5Hash: md5Hash(file),
		}
	}
}

func GetFileComparison(relativeFile string, homeFile string, projectFile string) FileComparison {
	hfile := getFileSituation(homeFile)
	pfile := getFileSituation(projectFile)
	inclination := NONE
	if !hfile.exists && pfile.exists {
		inclination = WILL_STORE
	} else if hfile.exists && !pfile.exists {
		inclination = WILL_RESTORE
	} else if !hfile.exists && !pfile.exists {
		inclination = BOTH_MISSING
	} else if hfile.exists && pfile.exists {
		// both exist, dig deeper to infer inclination
		if hfile.md5Hash == pfile.md5Hash {
			inclination = NONE
		} else if hfile.modTime.Before(pfile.modTime) {
			inclination = WILL_STORE
		} else if hfile.modTime.After(pfile.modTime) {
			inclination = WILL_RESTORE
		}
	}
	return FileComparison{
		RelativeFile: relativeFile,
		HomeFile:     homeFile,
		ProjectFile:  projectFile,
		Inclination:  inclination,
	}
}

func CompareProjectFiles(sideloadConfig config.SideloadConfig) ProjectFilesComparison {
	fileList := sideloadConfig.ProjectConfig.Files.Track
	sort.Strings(fileList)
	allComparisons := make([]FileComparison, len(fileList))
	for i, file := range fileList {
		homeFile := filepath.Join(sideloadConfig.HomeProjectDir, file)
		projectFile := filepath.Join(sideloadConfig.ProjectConfig.ProjectDir, file)
		comparison := GetFileComparison(file, homeFile, projectFile)
		allComparisons[i] = comparison
	}
	wouldRestore := FilterOfFileComparison(allComparisons, func(fc FileComparison) bool {
		return fc.Inclination == WILL_RESTORE
	})
	wouldStore := FilterOfFileComparison(allComparisons, func(fc FileComparison) bool {
		return fc.Inclination == WILL_STORE
	})
	wouldNothing := FilterOfFileComparison(allComparisons, func(fc FileComparison) bool {
		return fc.Inclination == NONE
	})
	bothMissing := FilterOfFileComparison(allComparisons, func(fc FileComparison) bool {
		return fc.Inclination == BOTH_MISSING
	})
	return ProjectFilesComparison{
		AllComparisons: allComparisons,
		WouldRestore:   wouldRestore,
		WouldStore:     wouldStore,
		WouldNothing:   wouldNothing,
		BothMissing:    bothMissing,
	}
}

func FilterOfFileComparison(collection []FileComparison, keep func(FileComparison) bool) []FileComparison {
	newCollection := make([]FileComparison, len(collection))
	j := 0
	for _, v := range collection {
		if keep(v) {
			newCollection[j] = v
			j++
		}
	}
	return newCollection[:j]
}
