package main

import (
	"context"
	"os/exec"
	"strings"
)

func init() {
	key := commandFunc{"!be offensive", "", "", categoryHidden}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if strings.HasPrefix(command, key.command) {
			out, err := exec.Command("fortune", "-o").Output()
			if err != nil {
				return "", err
			}
			return string(out), nil
		}
		return "", nil
	}
}
