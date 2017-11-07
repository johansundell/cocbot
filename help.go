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
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
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
			msg := "**COCBOT COMMANDS**\n"
			var c category
			for _, v := range keys {
				//fmt.Println(v.category)
				if v.category != categoryHidden {
					if c != v.category {
						msg += fmt.Sprintf("\n%s\n", v.category)
						c = v.category
					}
					msg += fmt.Sprintf("%s - %s\n", v.command, v.helpText)
				}
			}
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
			//sendEmbed(ch.ID, s, msg)
			//s.ChannelMessageSend(ch.ID, msg)
			//msg = "Sent you the help over a private channel, don't tell anyone our secret"
			return msg, nil
		}
		return "", nil
	}
}
