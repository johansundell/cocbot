package main

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/johansundell/cocapi"
)

var queryInsertUpdateMember = `INSERT INTO members (tag, name, created, last_updated, active) VALUES (?, ?, null, null, 1) ON DUPLICATE KEY UPDATE member_id=LAST_INSERT_ID(member_id), last_updated = NOW(), active = 1`
var channels []string
var isCocUnderUpdate bool
var failedTries int
var emailTo, emailFrom string

func reporter(s *discordgo.Session) {
	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				checkMembers(s)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func checkMembers(s *discordgo.Session) {
	members, err := cocClient.GetMembers(myClanTag)
	if err != nil {
		reportError(s, err)
		return
	}
	if isCocUnderUpdate {
		isCocUnderUpdate = false
		sendMessage(s, "Clash of Clans servers are up again")
	}
	failedTries = 0

	for _, m := range members.Items {
		if result, err := db.Exec(queryInsertUpdateMember, m.Tag, m.Name); err != nil {
			log.Println(err)
		} else {
			if id, err := result.LastInsertId(); err != nil {
				log.Println(err)
			} else {
				donations := 0
				if err := db.QueryRow("SELECT current_donations FROM members WHERE member_id = ?", id).Scan(&donations); err != nil {
					log.Println(err)
				} else {
					//log.Println(m.Donations, donations)
					if m.Donations != donations {
						if _, err := db.Exec("UPDATE members SET prev_donations = ?, current_donations = ?, last_donation_time = NOW() WHERE member_id = ?", donations, m.Donations, id); err != nil {
							log.Println(err)
						}
						if m.Donations > donations {
							if _, err := db.Exec("INSERT donations (member_id, ts, current_donations, prev_donations) VALUES (?, NOW(), ?, ?)", id, m.Donations, donations); err != nil {
								log.Println(err)
							}
						}
					}
				}
			}
		}
		if m.Role == "member" && m.Donations >= 1000 {
			log.Println("Found member that should be upgraded", m.Name)
			var alerted int
			db.QueryRow("SELECT alert_sent_donations FROM members WHERE tag = ?", m.Tag).Scan(&alerted)
			if alerted == 0 {
				//sendEmail("Member "+m.Name+" should be upgraded", "Member "+m.Name+" should be upgraded")
				sendMessage(s, "Member "+m.Name+" should be upgraded")
				db.Exec("UPDATE members SET alert_sent_donations = 1 WHERE tag = ?", m.Tag)
			}
		}
	}
}

func sendMessage(s *discordgo.Session, message string) {
	if len(channels) > 0 {
		for _, v := range channels {
			s.ChannelMessageSend(v, message)
		}
	}
}

func reportError(s *discordgo.Session, err error) {
	switch t := err.(type) {
	case *cocapi.ServerError:
		if t.ErrorCode == 503 {
			failedTries++
			if failedTries > 3 {
				if !isCocUnderUpdate {
					isCocUnderUpdate = true
					sendMessage(s, "Clash of Clans servers are down")
				}
			}
		}
		break
	default:
		log.Println("Fatal error coc:", t)
		break
	}
}
