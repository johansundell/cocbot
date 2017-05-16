package main

import "context"

func init() {
	key := commandFunc{"!show me the code", "To view me", "", categoryFun}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if key.command == command {
			return "You can find me here https://github.com/johansundell/cocbot", nil
		}
		return "", nil
	}
}
