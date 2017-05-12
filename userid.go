package main

import (
	"context"
)

func init() {
	key := commandFunc{"!userid", "Gives out your personal discord id", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if command == key.command {
			s, msg, err := getSessionsAndMessageFromContext(ctx)
			if err != nil {
				return "", err
			}
			userChannel, _ := s.UserChannelCreate(msg.Author.ID)
			s.ChannelMessageSend(userChannel.ID, "Your personal discord id is "+msg.Author.ID)
			return "Your personal id have been sent over a private channel", nil
		}
		return "", nil
	}
}
