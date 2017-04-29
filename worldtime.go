package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func init() {
	key := commandFunc{"!world time", "Show some of our members time", ""}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string) (string, error) {
		if strings.Contains(command, key.command) {
			m := make(map[string]string)
			m["Asia/Calcutta"] = "Brown and K2"
			m["Europe/Stockholm"] = "Sudde"
			m["Asia/Jakarta"] = "Tommy"
			m["EST"] = "James and Blac"

			msg := ""
			t := time.Now()
			for k, v := range m {
				l, err := time.LoadLocation(k)
				if err != nil {
					log.Println(err)
					continue
				}
				msg += fmt.Sprintf("The time for %s is %s\n", v, t.In(l).Format("15:04:05"))
			}
			return msg, nil
		}
		return "", nil
	}
}
