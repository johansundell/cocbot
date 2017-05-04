package main

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/johansundell/cocapi"
)

func init() {
	key := commandFunc{"!top war players", "To see our top war whores", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if key.command == command {
			list, err := cocClient.GetMembers(myClanTag)
			if err != nil {
				return "", err
			}
			players := []cocapi.Player{}
			for _, v := range list.Items {
				p, _ := cocClient.GetPlayerInfo(v.Tag)
				players = append(players, p)
			}
			sort.Slice(players, func(i, j int) bool { return players[i].WarStars > players[j].WarStars })
			msg := ""
			for i := 0; i < len(players); i++ {
				p := players[i]
				if i > 5 {
					break
				}
				msg += fmt.Sprintf("%d stars for %s\n", p.WarStars, p.Name)
			}
			return msg, nil
		}
		return "", nil
	}

	key2 := commandFunc{"!top [0-9]+ war players", "To see our top NN war whores", "[0-9]+"}
	botFuncs[key2] = func(command string, ctx context.Context) (string, error) {
		if found, _ := regexp.MatchString(key2.command, command); found {
			reg, err := regexp.Compile(key2.extracter)
			if err != nil {
				return "", err
			}
			result := reg.FindString(command)
			if result != "" {
				n, err := strconv.Atoi(result)
				if err != nil {
					return "", err
				}
				list, err := cocClient.GetMembers(myClanTag)
				if err != nil {
					return "", err
				}
				players := []cocapi.Player{}
				for _, v := range list.Items {
					p, _ := cocClient.GetPlayerInfo(v.Tag)
					players = append(players, p)
				}
				sort.Slice(players, func(i, j int) bool { return players[i].WarStars > players[j].WarStars })
				msg := ""
				for i := 0; i < len(players); i++ {
					p := players[i]
					if i > n-1 {
						break
					}
					msg += fmt.Sprintf("%d stars for %s\n", p.WarStars, p.Name)
				}
				return msg, nil
			}
		}
		return "", nil
	}
}
