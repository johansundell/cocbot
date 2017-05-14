package main

import "context"

func init() {
	key := commandFunc{"!status", "", "", categoryHidden}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if command == key.command {
			s, _, err := getSessionsAndMessageFromContext(ctx)
			if err != nil {
				return "", err
			}
			sendMessage(s, "Works", true)
			//return "Hi master", nil
		}
		return "", nil
	}
}
