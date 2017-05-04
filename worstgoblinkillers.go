package main

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/johansundell/cocapi"
)

type goblinPlayer struct {
	cocapi.Player
	goblinValue int
}

func init() {
	key := commandFunc{"!worst goblin killers", "To see our goblin hunters", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if key.command == command {
			list, err := cocClient.GetMembers(myClanTag)
			if err != nil {
				log.Println(err)
			}
			players := []goblinPlayer{}
			for _, v := range list.Items {
				p, _ := cocClient.GetPlayerInfo(v.Tag)
				for _, a := range p.Achievements {
					if a.Name == "Get those Goblins!" && a.Value < 150 {
						g := goblinPlayer{p, a.Value}
						players = append(players, g)
					}
				}
			}
			sort.Slice(players, func(i, j int) bool { return players[i].goblinValue < players[j].goblinValue })
			msg := ""
			for _, v := range players {
				msg += fmt.Sprintf("%s has %d Goblin stars left\n", v.Name, 150-v.goblinValue)
			}
			return msg, nil
		}
		return "", nil
	}
}
