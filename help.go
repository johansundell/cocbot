package main

import (
	"context"
	"fmt"
	"sort"
)

func init() {
	key := commandFunc{"!help", "List the help", "", catgoryHelp}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if command == key.command {
			var keys = make([]commandFunc, 0, len(botFuncs))
			for k := range botFuncs {
				keys = append(keys, k)
			}
			sort.Slice(keys, func(i, j int) bool {
				if keys[i].category == keys[j].category {
					return keys[i].command < keys[j].command
				}
				return keys[i].category < keys[j].category
			})
			msg := "**COCBOT COMMANDS**\n```"
			var c category
			for _, v := range keys {
				if v.category != categoryHidden {
					if c != v.category {
						msg += fmt.Sprintf("\n%s\n", v.category)
						c = v.category
					}
					msg += fmt.Sprintf("%s - %s\n", v.command, v.helpText)
				}
			}
			msg += "```"
			s, m, err := getSessionsAndMessageFromContext(ctx)
			if err != nil {
				return "", err
			}
			ch, _ := s.UserChannelCreate(m.Author.ID)
			u, _ := s.User("@me")
			pinned, _ := s.ChannelMessagesPinned(ch.ID)
			for _, v := range pinned {
				if v.Author.ID == u.ID {
					s.ChannelMessageUnpin(ch.ID, v.ID)
				}
			}
			mess, _ := s.ChannelMessageSend(ch.ID, msg)
			s.ChannelMessagePin(ch.ID, mess.ID)
			msg = "Sent you the help over a private channel"
			return msg, nil
		}
		return "", nil
	}
}
