package main

import (
	"context"
	"os/exec"
)

func init() {
	key := commandFunc{"!uptime", "To see how I am doing", "", categoryFun}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if command == key.command {
			out, err := exec.Command("uptime").Output()
			if err != nil {
				return "", err
			}
			msg := string(out)
			return msg, nil
		}
		return "", nil
	}
}
