package main

import (
	"context"
	"os/exec"
	"strings"
)

func init() {
	key := commandFunc{"!why do i suck at", "", "", categoryHidden}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
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
