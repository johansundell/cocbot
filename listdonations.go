package main

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"
)

func init() {
	key := commandFunc{"!list donations [0-9]+h [0-9]+m", "To list donation history, both hour and minutes must be given", "", categoryStats}
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		if found, _ := regexp.MatchString(key.command, command); found {
			fmt.Println("here")
			hours, err := getNumbersFromRegexp(command, "[0-9]+h")
			if err != nil {
				return "", err
			}
			min, err := getNumbersFromRegexp(command, "[0-9]+m")
			if err != nil {
				return "", err
			}
			t := time.Now()
			t = t.Add(-time.Hour * time.Duration(hours))
			t = t.Add(-time.Minute * time.Duration(min))
			don, err := getDonationsByTime(t, 10)
			if err != nil {
				return "", err
			}
			sort.Slice(don, func(i, j int) bool { return don[i].min > don[j].min })
			msg := fmt.Sprintf("Listing donations done after %d hours and %d minutes ago\n", hours, min)
			for _, v := range don {
				msg += v.String()
			}
			return msg, nil
		}
		return "", nil
	}
}

func getNumbersFromRegexp(s, pattern string) (int, error) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return 0, err
	}
	result := reg.FindString(s)
	result = result[:len(result)-1]
	i, err := strconv.Atoi(result)
	if err != nil {
		return 0, nil
	}
	return i, nil
}

var queryListDonationsTimeStamp = `
SELECT 
    (d.current_donations - d.prev_donations) AS diff,
    m.name,
    TIME_TO_SEC(TIMEDIFF(NOW(), d.ts)) AS since
FROM
    donations d
        JOIN
    members m ON m.member_id = d.member_id
WHERE d.ts >= ?
#ORDER BY since DESC
LIMIT 0 , ?
`

func getDonationsByTime(ts time.Time, numToFetch int) ([]donations, error) {
	rows, err := db.Query(queryListDonationsTimeStamp, ts.Format("2006-01-02 15:04:05"), numToFetch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	don := []donations{}
	for rows.Next() {
		d := donations{}
		rows.Scan(&d.amount, &d.name, &d.min)
		don = append(don, d)
	}
	return don, nil
}
