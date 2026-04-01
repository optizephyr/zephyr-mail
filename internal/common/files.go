package common

import (
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func EnsureDir(path string) error {
	if FileExists(path) {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

func Getwd() (string, error) {
	return os.Getwd()
}

func ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	wd, _ := os.Getwd()
	return filepath.Join(wd, path)
}

func WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
