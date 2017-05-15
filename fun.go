package main

import (
	"context"
	"fmt"
	"time"
)

func init() {
	key := commandFunc{"!game", "Lets play a game", "", categoryFun}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if command == key.command {
			s, m, err := getSessionsAndMessageFromContext(ctx)
			if err != nil {
				return "", err
			}
			n := 100
			str := "Formatting main hard drive %d%% left"
			msg, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(str, n))
			if err != nil {
				return "", err
			}
			ticker := time.NewTicker(1 * time.Second)
			quit := make(chan struct{})
			go func() {
				for {
					select {
					case <-ticker.C:
						if n == 0 {
							ticker.Stop()
							s.ChannelMessageDelete(m.ChannelID, msg.ID)
						}
						n = n - 5
						msg, err = s.ChannelMessageEdit(m.ChannelID, msg.ID, fmt.Sprintf(str, n))
						if err != nil {
							return
						}
					case <-quit:
						ticker.Stop()
						return
					}
				}
			}()

		}
		return "", nil
	}
}
