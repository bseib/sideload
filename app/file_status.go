package app

type Inclination int

const (
	NONE         Inclination = iota
	WILL_STORE   Inclination = iota
	WILL_RESTORE Inclination = iota
)

type FileStatus struct {
	homeFile    string
	projectFile string
	inclination Inclination
}

func GetFileStatus(homeFile string, projectFile string) FileStatus {
	return FileStatus{
		homeFile: homeFile,
		projectFile: projectFile,
		inclination: NONE,
	}
}