package main

import (
	"os/exec"
)

func init() {
	key := commandFunc{"!fetch your code", "", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if key.command == command {
			out, err := exec.Command("git", "pull").Output()
			if err != nil {
				return "", err
			}
			return string(out), err
		}
		return "", nil
	}
}
