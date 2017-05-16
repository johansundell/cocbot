package main

import (
	"context"
	"fmt"
	"sort"
	"time"
)

func init() {
	key := commandFunc{"!top donators", "To see our best donatots", "", categoryStats}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
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

func firstMonday(year int, month time.Month) int {
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	return (8-int(t.Weekday()))%7 + 1
}
