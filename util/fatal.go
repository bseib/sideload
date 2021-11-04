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

func Fatal(message string) {
	fmt.Println(message)
	os.Exit(1)
}