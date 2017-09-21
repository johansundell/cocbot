package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func reporterDuplicate(s *discordgo.Session) {
	ticker := time.NewTicker(15 * time.Minute)
	//ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				findDuplicate(s)

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func findDuplicate(s *discordgo.Session) {
	rows, err := db.Query("SELECT GROUP_CONCAT(name) as usernames, tag, count(*) AS c FROM members GROUP BY tag HAVING c > 1")

	if err != nil {
		log.Println(err)
	} else {
		msg := ""
		for rows.Next() {
			name, tag, c := "", "", 0
			rows.Scan(&name, &tag, &c)
			if clan, err := cocClient.GetClanInfo(myClanTag); err == nil {
				for _, m := range clan.MemberList {
					if m.Tag == tag {
						msg += fmt.Sprintf("Found duplicate member in clan %s, with tag %s has %d accounts", name, tag, c)
					}
				}
			}
		}
		rows.Close()
		if len(msg) > 0 {
			sendMessage(s, msg, true)
			//log.Println(msg)
		}
		//log.Println("test", msg)
	}
}
