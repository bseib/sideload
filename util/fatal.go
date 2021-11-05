package util

import (
	"fmt"
	"os"
)

func CheckFatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Fatal(v ...interface{}) {
	fmt.Println(fmt.Sprint(v...))
	os.Exit(1)
}
