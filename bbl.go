package main

import (
	"strings"
)

func init() {
	key := commandFunc{"bbl", "", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if strings.Contains(command, key.command) {
			return "No don't leave me here alone with ClanBot", nil
		}
		return "", nil
	}
}
