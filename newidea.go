package main

import (
	"context"
	"strings"
)

func init() {
	key := commandFunc{"!new idea", "", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if strings.HasPrefix(command, key.command) {
			return "That was a good idea master, I will write it down", nil
		}
		return "", nil
	}
}
