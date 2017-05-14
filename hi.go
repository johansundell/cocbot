package main

import (
	"context"
)

func init() {
	key := commandFunc{"!test", "Just a test functions", "", categoryFun}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if command == "!test" {
			return "Hi master", nil
		}
		return "", nil
	}
}
