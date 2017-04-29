package main

import (
	"strings"
)

func init() {
	key := commandFunc{"thanks cocbot", "", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if strings.Contains(command, key.command) {
			return "I live to serve you master", nil
		}
		return "", nil
	}
}
