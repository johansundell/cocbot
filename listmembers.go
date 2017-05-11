package main

import (
	"context"
	"strings"
)

type MemberList struct {
	Members []Member
}

type Member struct {
	Name   string
	Tag    string
	Active int
}

func init() {
	key := commandFunc{"!list members [name]", "To see current members", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		cmd := strings.Replace(key.command, " [name]", "", -1)
		if strings.HasPrefix(command, cmd) {
			mb, err := getMembers(command[len(cmd):])
			if err != nil {
				return "", err
			}
			msg := ""
			for _, v := range mb.Members {
				status := "inactive"
				if v.Active == 1 {
					status = "active"
				}
				msg += v.Name + " - " + v.Tag + " and status is " + status + "\n"
			}
			if len(msg) >= 2000 {
				msg = "Too broad search, limit it down"
			}
			return msg, nil
		}
		return "", nil
	}
}

func getMembers(search string) (MemberList, error) {
	search = strings.TrimSpace(search)
	rows, err := db.Query("SELECT tag, name, active FROM members WHERE name LIKE ? ORDER BY active, name", "%"+search+"%")
	if err != nil {
		return MemberList{}, err
	}
	defer rows.Close()
	mb := MemberList{}
	for rows.Next() {
		m := Member{}
		rows.Scan(&m.Tag, &m.Name, &m.Active)
		mb.Members = append(mb.Members, m)
	}
	return mb, nil
}
