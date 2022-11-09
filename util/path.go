package util

import (
	"os"
	"path/filepath"
)

func CreateFolder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return os.MkdirAll(folder, os.ModePerm)
	}
	return nil
}

func EnsureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModeDir)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func DeletedFolder(folder string) error {
	if _, err := os.Stat(folder); !os.IsNotExist(err) {
		return os.RemoveAll(folder)
	}
	return nil
}

func GetFolderRoot() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	return dir, nil
}
