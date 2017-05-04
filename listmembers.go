package main

import (
	"context"
	"log"
	"strings"
)

type MemberList struct {
	Members []Member
}

type Member struct {
	Name string
	Tag  string
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
				msg += v.Name + " - " + v.Tag + "\n"
			}
			return msg, nil
		}
		return "", nil
	}
}

func getMembers(search string) (MemberList, error) {
	search = strings.TrimSpace(search)
	log.Println(search)
	rows, err := db.Query("SELECT tag, name FROM members WHERE active = 1 AND name LIKE ? ORDER BY name", "%"+search+"%")
	if err != nil {
		return MemberList{}, err
	}

	mb := MemberList{}
	for rows.Next() {
		m := Member{}
		rows.Scan(&m.Tag, &m.Name)
		mb.Members = append(mb.Members, m)
	}
	return mb, nil
}
