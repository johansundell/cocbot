package main

import (
	"context"
	"strings"
)

func init() {
	key := commandFunc{"bbl", "", "", categoryHidden}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if strings.Contains(command, key.command) {
			return "No don't leave me here alone", nil
		}
		return "", nil
	}
}
