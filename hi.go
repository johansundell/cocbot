package main

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func init() {
	key := commandFunc{"!test", "Just a test functions", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if command == "!test" {
			if s := ctx.Value("sess"); s != nil {
				if m := ctx.Value("msg"); m != nil {
					s.(*discordgo.Session).ChannelMessageSend(m.(*discordgo.MessageCreate).ChannelID, "testing context")
				}
			}
			return "Hi master", nil
		}
		return "", nil
	}
}
