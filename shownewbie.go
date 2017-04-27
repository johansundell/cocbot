package main

import (
	"fmt"
	"sort"

	"github.com/johansundell/cocapi"
)

func init() {
	key := commandFunc{"!show newbie", "To see our newest members and their donations", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if key.command == command {
			list, _ := cocClient.GetMembers(myClanTag)
			players := []cocapi.Player{}
			for _, v := range list.Items {
				if v.Role == "member" {
					if p, err := cocClient.GetPlayerInfo(v.Tag); err == nil {
						players = append(players, p)
					}
				}
			}
			sort.Slice(players, func(i, j int) bool { return players[i].Donations > players[j].Donations })
			msg := ""
			for _, v := range players {
				msg += fmt.Sprintf("%d %s th%d total hero lvl %d\n", v.Donations, v.Name, v.TownHallLevel, getHeroLvl(v))
			}
			return msg, nil
		}
		return "", nil
	}
}

func getHeroLvl(p cocapi.Player) int {
	tot := 0
	for _, v := range p.Heroes {
		tot += v.Level
	}
	return tot
}
