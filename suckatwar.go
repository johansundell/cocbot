package main

import (
	"os/exec"
	"strings"
)

func init() {
	key := commandFunc{"!why do i suck at", "", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if strings.HasPrefix(command, key.command) {
			out, err := exec.Command("fortune", "tao").Output()
			if err != nil {
				return "", err
			}
			return string(out), nil
		}
		return "", nil
	}
}
