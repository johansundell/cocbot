package main

import (
	"context"
	"fmt"
	"time"
)

type receive struct {
	donations
}

func (v receive) String() string {
	return fmt.Sprintf("%s got %d troops %s ago\n", v.name, v.amount, (time.Duration(v.min) * time.Second).String())
}

func init() {
	key := commandFunc{"!what to name 2", "To list the last 10 to receive troops", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if key.command == command {
			don, err := getReceive(10)
			if err != nil {
				return "", err
			}
			msg := ""
			for _, v := range don {
				msg += fmt.Sprintf("%s got %d troops %s ago\n", v.name, v.amount, (time.Duration(v.min) * time.Second).String())
			}
			return msg, nil
		}
		return "", nil
	}
}

var queryRecive = `
SELECT 
    (r.current - r.prev) AS diff,
    m.name,
    TIME_TO_SEC(TIMEDIFF(NOW(), r.ts)) AS since
FROM
    receive r
        JOIN
    members m ON m.member_id = r.member_id
ORDER BY r.receive_id DESC
LIMIT 0 , ?
`

func getReceive(numToFetch int) ([]receive, error) {
	rows, err := db.Query(queryRecive, numToFetch)
	if err != nil {
		return nil, err
	}
	don := []receive{}
	for rows.Next() {
		d := receive{}
		rows.Scan(&d.amount, &d.name, &d.min)
		don = append(don, d)
	}
	return don, nil
}
