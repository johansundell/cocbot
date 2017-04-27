package main

func init() {
	key := commandFunc{"!send me nude pics", "To see me nude", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if key.command == command {
		}
		return "", nil
	}
}
