package main

import "os/exec"

func init() {
	key := commandFunc{"!why do I suck at war", "Hopefully somwthing wise", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if key.command == command {
			out, err := exec.Command("fortune", "tao").Output()
			if err != nil {
				return "", err
			}
			return string(out), nil
		}
		return "", nil
	}
}
