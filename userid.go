package main

import (
	"context"
)

func init() {
	key := commandFunc{"!sudde", "Gives out your personal discord id", "", categoryAdmin}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if command == key.command {
			s, msg, err := getSessionsAndMessageFromContext(ctx)
			if err != nil {
				return "", err
			}
			userChannel, _ := s.UserChannelCreate(msg.Author.ID)
			m, _ := s.ChannelMessageSend(userChannel.ID, "Your personal discord id is "+msg.Author.ID)
			s.ChannelMessagePin(userChannel.ID, m.ID)
			return "Your personal id have been sent over a private channel", nil
		}
		return "", nil
	}
}
