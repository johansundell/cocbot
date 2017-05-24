package main

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func init() {
	key := commandFunc{"!test", "Just a test functions", "", categoryFun}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if command == key.command {
			s, m, err := getSessionsAndMessageFromContext(ctx)
			if err != nil {
				return "", err
			}
			sendMessage(s, "<@&"+everyoneId+">"+"test message, sorry\nTesting for blac", false)
			em := discordgo.MessageEmbed{}
			em.Title = "test\nhej hopp och lite mera text"

			em.Footer = &discordgo.MessageEmbedFooter{}
			em.Footer.Text = "cocbot"

			em.Author = &discordgo.MessageEmbedAuthor{}
			em.Author.Name = "sudde"
			em.Author.IconURL = "https://avatars1.githubusercontent.com/u/6422482?v=3&s=400"
			em.Author.URL = "http://pixpro.net"

			em.Color = 11584734

			em.Description = "en liten text som testar detta\noch har stöd för radbrytning"

			fields := make([]*discordgo.MessageEmbedField, 0)
			f := discordgo.MessageEmbedField{}
			f.Name = "th"
			f.Value = "9"
			f.Inline = true
			fields = append(fields, &f)

			f2 := discordgo.MessageEmbedField{}
			f2.Name = "stars"
			f2.Value = "3"
			f2.Inline = true
			fields = append(fields, &f2)

			em.Fields = fields

			_, err = s.ChannelMessageSendEmbed(m.ChannelID, &em)
			if err != nil {
				return "", err
			}

			//return "Hi master", nil
		}
		return "", nil
	}
}
