package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func reporterDuplicate(s *discordgo.Session) {
	ticker := time.NewTicker(15 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				rows, err := db.Query("SELECT GROUP_CONCAT(name) as usernames, tag, count(*) AS c FROM members GROUP BY tag HAVING c > 1")

				if err != nil {
					log.Println(err)
				} else {
					msg := ""
					for rows.Next() {
						name, tag, c := "", "", 0
						rows.Scan(&name, &tag, &c)
						msg += fmt.Sprintf("Found duplicate member in clan %s, with tag %s has %d accounts", name, tag, c)
					}
					rows.Close()
					if len(msg) > 0 {
						sendMessage(s, msg)
					}
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
