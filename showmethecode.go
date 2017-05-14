package main

import "context"

func init() {
	key := commandFunc{"!show me the code", "To view me", "", categoryFun}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if key.command == command {
			return "You can find me here https://github.com/johansundell/cocbot", nil
		}
		return "", nil
	}
}
