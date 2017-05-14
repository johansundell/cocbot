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
	botFuncs[keyTopDonators] = func(command string, ctx context.Context) (string, error) {
		if command == keyTopDonators.command {
			rows, err := db.Query(sqlQueryTopTodayDonators, 10)
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
			msg := "Todays top donators are, reset at " + f(getDuration().String()) + "\n"
			for _, v := range result {
				msg += fmt.Sprintf("%d troops by %s\n", v.amount, v.name)
			}
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
        AND d.ts < CURDATE() + INTERVAL 1 DAY AND m.active = 1
GROUP BY m.member_id
ORDER BY diff DESC
LIMIT 0 , ?
`
