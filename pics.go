package main

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

func init() {
	key := commandFunc{"!send me nude pics", "To see me nude", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if key.command == command {
		}
		return "", nil
	}
}

func getRandomImage() (string, error) {
	m := make(map[int]string)
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := path.Dir(ex)
	files, err := filepath.Glob(exPath + string(os.PathSeparator) + "*.jpg")
	if err != nil {
		return "", err
	}
	for k, v := range files {
		m[k] = v
	}
	for _, v := range m {
		return v, nil
	}
	return "", errors.New("No files found")
}
