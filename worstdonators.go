package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/johansundell/cocapi"
)

func init() {
	key := commandFunc{"!worst donators", "To see our worst donators", "", categoryStats}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if key.command == command {
			list, err := cocClient.GetClanInfo(myClanTag)
			if err != nil {
				return "", err
			}

			sort.Sort(sort.Reverse(cocapi.DonationRatio(list.MemberList)))

			msg := ""
			for i := 0; i < len(list.MemberList) && i < 15; i++ {
				v := list.MemberList[i]
				msg += fmt.Sprintf("%s has donation difference %s troops\n", v.Name, humanize.FormatInteger("### ###,", v.Donations-v.DonationsReceived))
			}
			fmt.Println(msg)
			return msg, nil
		}
		return "", nil
	}
}
