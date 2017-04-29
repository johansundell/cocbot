package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func getTime() time.Time {
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 23, 58, 0, 0, time.Local)
	return t
}

func getDuration() time.Duration {
	t := getTime()
	return t.Sub(time.Now())
}

func test(s *discordgo.Session) {
	//now := time.Now()
	t := getTime()
	//t = t.Add(time.Minute * time.Duration(2))
	//fmt.Println(t)
	//return
	d := t.Sub(time.Now())
	fmt.Println(d)
	//ticker := time.NewTicker(1 * time.Minute)
	ticker := time.NewTimer(d)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("ticker hit")
				//lockMap.RLock()
				msg, err := botFuncs[keyTopDonators]("!top donators today")
				//lockMap.RUnlock()
				if err != nil {
					log.Println(err)
				} else {
					fmt.Println(msg)
					sendMessage(s, msg)
				}
				t = getTime()
				t = t.Add(time.Hour * time.Duration(24))
				d = t.Sub(time.Now())
				ticker.Reset(d)
				fmt.Println("ran", d)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
