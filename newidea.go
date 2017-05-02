package main

import (
	"strings"
)

func init() {
	key := commandFunc{"!new idea", "", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if strings.HasPrefix(command, key.command) {
			return "That was a good idea master, I will write it down", nil
		}
		return "", nil
	}
}
