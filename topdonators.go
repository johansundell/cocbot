package main

import (
	"fmt"
	"sort"
)

func init() {
	key := commandFunc{"!top donators", "To see our best donatots", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if key.command == command {
			list, _ := cocClient.GetMembers(myClanTag)
			sort.Slice(list.Items, func(i, j int) bool { return list.Items[i].Donations > list.Items[j].Donations })
			msg := ""
			for n := 0; n < len(list.Items); n++ {
				v := list.Items[n]
				msg += fmt.Sprintf("%d %s\n", v.Donations, v.Name)
				if n > 8 {
					break
				}
			}
			return msg, nil
		}
		return "", nil
	}
}
