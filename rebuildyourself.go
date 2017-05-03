package main

import "os/exec"

func init() {
	key := commandFunc{"!rebuild yourself", "", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if key.command == command {
			out, err := exec.Command("go", "build").Output()
			if err != nil {
				return "", err
			}
			msg := ""
			if string(out) == "" {
				msg = "Done now master"
			} else {
				return string(out), nil
			}
			return msg, nil
			/*out, err = exec.Command("go", "install").Output()
			if err != nil {
				return "", err
			}
			if string(out) == "" {
				return msg, nil
			} else {
				return string(out), nil
			}*/
		}
		return "", nil
	}
}
