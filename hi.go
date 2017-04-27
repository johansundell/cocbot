package main

func init() {
	key := commandFunc{"!test", "Just a test functions", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if command == "!test" {
			return "Hi master", nil
		}
		return "", nil
	}
}
