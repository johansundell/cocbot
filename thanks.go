package main

import (
	"context"
	"strings"
)

func init() {
	key := commandFunc{"thanks cocbot", "", "", categoryHidden}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if strings.Contains(command, key.command) {
			return "I live to serve you master", nil
		}
		return "", nil
	}
}
