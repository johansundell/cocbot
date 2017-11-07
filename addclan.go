package main

import (
	"context"
	"log"
	"strings"
)

func init() {
	key := commandFunc{"!add clan [name]", "To add a new clan to be watched", "", categoryAdmin}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		cmd := strings.Replace(key.command, " [name]", "", -1)
		if strings.HasPrefix(command, cmd) {
			tag := strings.TrimSpace(command[len(cmd):])
			info, err := cocClient.GetClanInfo(tag)
			if err != nil {
				return "", err
			}
			var name string
			if err := db.QueryRow("SELECT name from clans WHERE tag = ?", info.Tag).Scan(&name); err != nil {
				log.Println(err)
				if _, err := db.Exec("INSERT INTO clans (tag, name) VALUES (?, ?)", info.Tag, info.Name); err != nil {
					return "", err
				} else {
					return info.Name + " has been added to the list", nil
				}
			} else {
				return info.Name + " was already on the list", nil
			}
		}
		return "", nil
	}
}
