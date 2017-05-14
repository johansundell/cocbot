package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/johansundell/cocapi"
)

type friendPlayer struct {
	cocapi.Player
	friendValue int
}

func init() {
	key := commandFunc{"!best friends", "To see our 10 best friends", "", categoryStats}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if key.command == command {
			list, err := cocClient.GetMembers(myClanTag)
			if err != nil {
				return "", err
			}
			players := []friendPlayer{}
			for _, v := range list.Items {
				p, _ := cocClient.GetPlayerInfo(v.Tag)
				for _, a := range p.Achievements {
					if a.Name == "Friend in Need" {
						g := friendPlayer{p, a.Value}
						players = append(players, g)
					}
				}
			}
			sort.Slice(players, func(i, j int) bool { return players[i].friendValue > players[j].friendValue })
			msg := ""
			for i := 0; i < len(players) && i < 10; i++ {
				v := players[i]
				msg += fmt.Sprintf("%s has donated %s troops\n", v.Name, humanize.FormatInteger("### ###,", v.friendValue))
			}
			return msg, nil
		}
		return "", nil
	}
}
