package main

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	key := commandFunc{"!clan hoppers", "List clan hoppers", "", categoryAdmin}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if command == key.command {
			s, msg, err := getSessionsAndMessageFromContext(ctx)
			if err != nil {
				return "", err
			}
			if doesMemberHasAdminAccess(s, msg) || isSudde(msg) {
				rows, err := db.Query("SELECT name FROM members WHERE active = 1 AND exited > 0")
				defer rows.Close()
				if err != nil {
					return "", err
				}
				members := make([]string, 0)
				for rows.Next() {
					var name string
					rows.Scan(&name)
					members = append(members, name)
				}
				sort.Strings(members)
				msg := "**These are our clan hoppers\n**"
				if len(members) == 0 {
					return "**No hoppers found\n**", nil
				}
				if len(members) > 0 {

					return msg + strings.Join(members, "\n"), nil
				}
			} else {
				return securityMessage, nil
			}
		}
		return "", nil
	}
}

func getSessionsAndMessageFromContext(ctx context.Context) (*discordgo.Session, *discordgo.MessageCreate, error) {
	if s, err := ctx.Value("sess").(*discordgo.Session); err {
		if msg, err := ctx.Value("msg").(*discordgo.MessageCreate); err {
			return s, msg, nil
		}
	}
	return nil, nil, errors.New("Could not get session and message from contect")
}

func doesMemberHasAdminAccess(s *discordgo.Session, msg *discordgo.MessageCreate) bool {
	member, _ := s.GuildMember(guild, msg.Author.ID)
	sort.Strings(member.Roles)
	i := sort.SearchStrings(member.Roles, coLeaderId)
	if i < len(member.Roles) && member.Roles[i] == coLeaderId {
		// Co-Leader
		return true
	}
	i = sort.SearchStrings(member.Roles, leaderId)
	if i < len(member.Roles) && member.Roles[i] == leaderId {
		// Leader
		return true
	}
	return false
}
