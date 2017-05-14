package main

import (
	"context"
	"fmt"
	"strings"
)

func init() {
	key := commandFunc{"!last donations to [name]", "Get the 10 last donations to that member", "", categoryStats}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		cmd := strings.Replace(key.command, " [name]", "", -1)
		if strings.HasPrefix(command, cmd) {
			fmt.Println("here")
			name := strings.TrimSpace(command[len(cmd):])
			don, err := getDonationsTo(name, 10)
			if err != nil {
				return "", err
			}
			msg := "This is the last 10 donations to " + name + "\n"
			for _, v := range don {
				msg += v.String()
			}
			return msg, nil
		}
		return "", nil
	}
}

var queryListDonationsTo = `
SELECT 
    (d.current - d.prev) AS diff,
    m.name,
    TIME_TO_SEC(TIMEDIFF(NOW(), d.ts)) AS since
FROM
    receive d
        JOIN
    members m ON m.member_id = d.member_id
WHERE m.name LIKE ?
ORDER BY d.receive_id DESC
LIMIT 0 , ?
`

func getDonationsTo(name string, numToFetch int) ([]receive, error) {
	rows, err := db.Query(queryListDonationsTo, name, numToFetch)
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
