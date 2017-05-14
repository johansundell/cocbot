package main

import (
	"context"
	"os/exec"
)

func init() {
	key := commandFunc{"!uptime", "To see how I am doing", "", categoryFun}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if command == key.command {
			out, err := exec.Command("go", "build").Output()
			if err != nil {
				return "", err
			}
			msg := string(out)
			return msg, nil
		}
		return "", nil
	}
}
