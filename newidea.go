package main

import (
	"context"
	"strings"
)

func init() {
	key := commandFunc{"!new idea", "", "", categoryHidden}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if strings.HasPrefix(command, key.command) {
			return "That was a good idea master, I will write it down", nil
		}
		return "", nil
	}
}
