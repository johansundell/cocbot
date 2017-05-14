package main

import (
	"context"
	"strings"
)

func init() {
	key := commandFunc{"!remove hopper [name]", "Remove hopper warnings", "", categoryAdmin}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(command string, ctx context.Context) (string, error) {
		cmd := strings.Replace(key.command, " [name]", "", -1)
		if strings.HasPrefix(command, cmd) {
			s, msg, err := getSessionsAndMessageFromContext(ctx)
			if err != nil {
				return "", err
			}
			if doesMemberHasAdminAccess(s, msg) || isSudde(msg) {
				name := strings.TrimSpace(command[len(cmd):])
				if len(name) > 0 {
					res, err := db.Exec("UPDATE members SET exited = 0 WHERE name = ?", name)
					if err != nil {
						return "", nil
					}
					if n, err := res.RowsAffected(); err == nil && n > 0 {
						return "Removed hopper warning for " + name, nil
					}
					return "Could not find member " + name + " or (s)he had no warnings", nil
				}
			} else {
				return securityMessage, nil
			}
		}
		return "", nil
	}
}
