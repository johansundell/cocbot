package main

import (
	"context"
	"os/exec"
)

func init() {
	key := commandFunc{"!fetch your code", "", "", categoryAdmin}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if key.command == command {
			s, m, err := getSessionsAndMessageFromContext(ctx)
			if err != nil {
				return "", err
			}
			s.ChannelMessageSend(m.ChannelID, "On it master")
			s.ChannelTyping(m.ChannelID)
			out, err := exec.Command("git", "pull").Output()
			if err != nil {
				return "", err
			}
			return string(out), err
		}
		return "", nil
	}
}
