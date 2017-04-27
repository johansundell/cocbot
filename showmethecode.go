package main

func init() {
	key := commandFunc{"!show me the code", "To view me", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if key.command == command {
			return "You can find me here https://github.com/johansundell/cocbot", nil
		}
		return "", nil
	}
}
