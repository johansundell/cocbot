package main

import (
	"fmt"
	"strings"
)

type donations struct {
	name   string
	min    int64
	amount int64
}

func init() {
	key := commandFunc{"!last donations", "To list the last 10 donations", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if key.command == command {
			don, err := getDonations(10)
			if err != nil {
				return "", err
			}
			msg := ""
			for _, v := range don {
				msg += fmt.Sprintf("%s donated %d troops %d minutes ago\n", v.name, v.amount, v.min)
			}
			return msg, nil
		}
		return "", nil
	}

	key = commandFunc{"!last donations for [name]", "To see the last donations by that member", ""}
	botFuncs[key] = func(command string) (string, error) {
		cmd := strings.Replace(key.command, " [name]", "", -1)
		if strings.HasPrefix(command, cmd) {
			name := command[len(cmd):]
			don, err := getUserDonations(name, 10)
			if err != nil {
				return "", err
			}
			msg := ""
			if len(don) != 0 {
				msg += "These are the last donations by " + name + "\n"
				for _, v := range don {
					msg += fmt.Sprintf("%d troops %d minutes ago\n", v.amount, v.min)
				}
			}
			return msg, nil
		}
		return "", nil
	}
}

var queryDonations = `
SELECT 
    (d.current_donations - d.prev_donations) AS diff,
    m.name,
    ROUND(TIME_TO_SEC(TIMEDIFF(NOW(), d.ts)) / 60) AS since
FROM
    donations d
        JOIN
    members m ON m.member_id = d.member_id
ORDER BY d.donate_id DESC
LIMIT 0 , ?
`

func getDonations(numToFetch int) ([]donations, error) {
	rows, err := db.Query(queryDonations, numToFetch)
	if err != nil {
		return nil, err
	}
	don := []donations{}
	for rows.Next() {
		d := donations{}
		rows.Scan(&d.amount, &d.name, &d.min)
		don = append(don, d)
	}
	return don, nil
}

var queryUserDonations = `
SELECT 
    (d.current_donations - d.prev_donations) AS diff,
    m.name,
    ROUND(TIME_TO_SEC(TIMEDIFF(NOW(), d.ts)) / 60) AS since
FROM
    donations d
        JOIN
    members m ON m.member_id = d.member_id
WHERE m.name LIKE ?
ORDER BY d.donate_id DESC
LIMIT 0 , ?
`

func getUserDonations(name string, numToFetch int) ([]donations, error) {
	rows, err := db.Query(queryUserDonations, name, numToFetch)
	if err != nil {
		return nil, err
	}
	don := []donations{}
	for rows.Next() {
		d := donations{}
		rows.Scan(&d.amount, &d.name, &d.min)
		don = append(don, d)
	}
	return don, nil
}
