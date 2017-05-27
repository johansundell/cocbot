package main

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func sniffer() {
	ticker := time.NewTicker(10 * time.Second)
	getMembersData(myClanTag)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := getMembersData(myClanTag); err != nil {
					log.Panicln(err)
				}
			}
		}
	}()
}

func getMembersData(clan string) error {
	members, err := cocClient.GetMembers(clan)
	if err != nil {
		//log.Println(err)
		//reportError(err)
		return err
	}

	if isCocUnderUpdate {
		isCocUnderUpdate = false
		//sendEmail("COC Alert", "Servers are up again")
	}
	failedTries = 0

	var ids = make([]string, 0)
	for _, m := range members.Items {
		if result, err := db.Exec(queryInsertUpdateMember, m.Tag, m.Name); err != nil {
			log.Println(err)
		} else {
			if id, err := result.LastInsertId(); err != nil {
				log.Println(err)
			} else {
				donations := 0
				if err := db.QueryRow("SELECT current_donations FROM members WHERE member_id = ?", id).Scan(&donations); err != nil {
					log.Println(err, "no 1")
				}
				//log.Println(m.Donations, donations)
				if m.Donations != donations {
					if _, err := db.Exec("UPDATE members SET prev_donations = ?, current_donations = ?, last_donation_time = NOW() WHERE member_id = ?", donations, m.Donations, id); err != nil {
						log.Println(err)
					}
					//if m.Donations > donations {
					if _, err := db.Exec("INSERT donations (member_id, ts, current_donations, prev_donations) VALUES (?, NOW(), ?, ?)", id, m.Donations, donations); err != nil {
						log.Println(err, "")
					}
					//}
				}

				received := 0
				if err := db.QueryRow("SELECT current_rec FROM members WHERE member_id = ?", id).Scan(&received); err != nil {
					log.Println(err)
				}
				if m.DonationsReceived != received {
					if _, err := db.Exec("UPDATE members SET current_rec = ? WHERE member_id = ?", m.DonationsReceived, id); err != nil {
						log.Println(err)
					}
					if _, err := db.Exec("INSERT receive (member_id, ts, current, prev) VALUES ( ?, NOW(), ?, ? )", id, m.DonationsReceived, received); err != nil {
						log.Println(err)
					}
				}
				ids = append(ids, strconv.Itoa(int(id)))
			}
		}
		if m.Role == "member" && m.Donations >= 1000 {
			log.Println("Found member that should be upgraded", m.Name)
			var alerted int
			db.QueryRow("SELECT alert_sent_donations FROM members WHERE tag = ?", m.Tag).Scan(&alerted)
			if alerted == 0 {
				//sendEmail("Member "+m.Name+" should be upgraded", "Member "+m.Name+" should be upgraded")
				db.Exec("UPDATE members SET alert_sent_donations = 1 WHERE tag = ?", m.Tag)
			}
		}
	}
	db.Exec("UPDATE members SET exited = NOW() WHERE member_id NOT IN (" + strings.Join(ids, ", ") + ") AND active = 1")
	db.Exec("UPDATE members SET active = 0 WHERE member_id NOT IN (" + strings.Join(ids, ", ") + ")")
	//log.Println("done members func")
	return nil
}
