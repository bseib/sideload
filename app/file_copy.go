package app

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src string, dest string) (bytesWritten int64, err error) {
	err = os.MkdirAll(filepath.Dir(dest), 0700)
	if err != nil {
		return 0, err
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()
	destFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()
	bytesWritten, err = io.Copy(destFile, srcFile)
	if err != nil {
		return bytesWritten, err
	}
	return bytesWritten, nil
}
