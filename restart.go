package main

import "os/exec"

func init() {
	key := commandFunc{"!restart yourself", "", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if key.command == command {
			out, err := exec.Command("/etc/init.d/cocbot", "restart").Output()
			if err != nil {
				return "", err
			}
			if string(out) == "" {
				return "Done now master", nil
			}
			return string(out), nil
		}
		return "", nil
	}
}
