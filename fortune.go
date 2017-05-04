package main

import (
	"context"
	"os/exec"
)

func init() {
	key := commandFunc{"!fortune", "To get a fortune cookie", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if key.command == command {
			out, err := exec.Command("fortune", "computers", "men-women").Output()
			if err != nil {
				return "", err
			}
			return string(out), nil
		}
		return "", nil
	}
}
