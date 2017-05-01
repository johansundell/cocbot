package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type donations struct {
	name   string
	min    int64
	amount int64
}

func (v donations) String() string {
	return fmt.Sprintf("%s donated %d troops %s ago\n", v.name, v.amount, (time.Duration(v.min) * time.Second).String())
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
				msg += v.String()
			}
			rec, err := getReceive(10)
			if err != nil {
				return "", err
			}
			msg += "--------\n"
			for _, v := range rec {
				msg += v.String()
			}
			log.Println(len(msg))
			return msg, nil
		}
		return "", nil
	}

	key2 := commandFunc{"!last donations for [name]", "To see the last donations by that member", ""}
	botFuncs[key2] = func(command string) (string, error) {
		cmd := strings.Replace(key2.command, " [name]", "", -1)
		if strings.HasPrefix(command, cmd) {
			name := strings.TrimSpace(command[len(cmd):])
			don, err := getUserDonations(name, 10)
			if err != nil {
				return "", err
			}
			msg := ""
			if len(don) != 0 {
				msg += "These are the last donations by " + name + "\n"
				for _, v := range don {
					msg += v.String()
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
    TIME_TO_SEC(TIMEDIFF(NOW(), d.ts)) AS since
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
    TIME_TO_SEC(TIMEDIFF(NOW(), d.ts)) AS since
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
