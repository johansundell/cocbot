package main

import (
	"context"
	"os/exec"
)

func init() {
	key := commandFunc{"!rebuild yourself", "", "", categoryAdmin}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
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
