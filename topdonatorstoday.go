package main

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type topDonator struct {
	name   string
	amount int64
}

var keyTopDonators = commandFunc{"!top donators today", "To see the our best donators today", "", categoryStats}

func init() {
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[keyTopDonators] = func(ctx context.Context, command string) (string, error) {
		if command == keyTopDonators.command {
			var msg string
			clans, err := db.Query("SELECT clan_id, name FROM clans")
			if err != nil {
				return "", err
			}
			for clans.Next() {
				var clanId int
				var clanName string
				if err := clans.Scan(&clanId, &clanName); err != nil {
					return "", err
				}
				rows, err := db.Query(sqlQueryTopTodayDonators, clanId, 10)
				if err != nil {
					return "", err
				}
				result := []topDonator{}
				for rows.Next() {
					r := topDonator{}
					if err := rows.Scan(&r.amount, &r.name); err != nil {
						log.Println(err)
					}
					result = append(result, r)
				}
				f := func(s string) string {
					n := strings.LastIndex(s, ".")
					return s[:n] + "s"
				}
				msg = "Todays top donators are in " + clanName + ", reset at " + f(getDuration().String()) + "\n"
				for _, v := range result {
					msg += fmt.Sprintf("%d troops by %s\n", v.amount, v.name)
				}
				s, _, err := getSessionsAndMessageFromContext(ctx)
				channel, found := ctx.Value("channel").(string)
				if found {
					log.Println("No channel")
					return "", nil
				}

				if err != nil {
					return "", err
				}

				s.ChannelMessageSend(channel, msg)
			}
			msg = ""
			return msg, nil
		}
		return "", nil
	}
}

var sqlQueryTopTodayDonators = `
SELECT 
    (SUM(d.current_donations) - SUM(d.prev_donations)) AS diff,
    m.name
FROM
    donations d
        JOIN
    members m ON m.member_id = d.member_id
WHERE
    d.ts >= CURDATE()
        AND d.ts < CURDATE() + INTERVAL 1 DAY AND m.active = 1 AND clan_id = ?
GROUP BY m.member_id
ORDER BY diff DESC
LIMIT 0 , ?
`
