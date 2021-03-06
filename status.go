package main

import (
	"context"
	"fmt"
)

func init() {
	key := commandFunc{"!status", "", "", categoryHidden}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if command == key.command {
			s, m, err := getSessionsAndMessageFromContext(ctx)
			if err != nil {
				return "", err
			}
			if isSudde(m) {
				if myChannel != "" {
					c, _ := s.Channel(myChannel)
					msg := fmt.Sprintf("Reporting on channel <#%s>", c.ID)
					sendMessage(s, msg, false)

				}
			} else {
				return securityMessage, nil
			}
			//return "Hi master", nil
		}
		return "", nil
	}
}
