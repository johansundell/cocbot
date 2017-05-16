package main

import (
	"context"
	"fmt"
	"log"
)

func init() {
	key := commandFunc{"!what to name", "To see the our best recivers today", "", categoryStats}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if command == key.command {
			rows, err := db.Query(sqlQueryTopReciveDonators, 10)
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
			msg := "Todays top recivers are\n"
			for _, v := range result {
				msg += fmt.Sprintf("%d troops to %s\n", v.amount, v.name)
			}
			return msg, nil
		}
		return "", nil
	}
}

var sqlQueryTopReciveDonators = `
SELECT 
    (SUM(d.current) - SUM(d.prev)) AS diff,
    m.name
FROM
    receive d
        JOIN
    members m ON m.member_id = d.member_id
WHERE
    d.ts >= CURDATE()
        AND d.ts < CURDATE() + INTERVAL 1 DAY
GROUP BY m.member_id
ORDER BY diff DESC
LIMIT 0 , ?
`
