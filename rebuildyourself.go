package main

import (
	"context"
	"os/exec"
)

func init() {
	key := commandFunc{"!rebuild yourself", "", "", categoryAdmin}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if key.command == command {
			s, m, err := getSessionsAndMessageFromContext(ctx)
			if err != nil {
				return "", err
			}
			sendEmbed(m.ChannelID, s, "Working on it master, new version of self soon done")
			s.ChannelTyping(m.ChannelID)
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
