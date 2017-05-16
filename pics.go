package main

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	key := commandFunc{"!send me nude pics", "To see me nude", "", categoryFun}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx context.Context, command string) (string, error) {
		if key.command == command {
			if s := ctx.Value("sess"); s != nil {
				if m := ctx.Value("msg"); m != nil {
					str, err := getRandomImage()
					if err != nil {
						return "", err
					}
					f, err := os.Open(str)
					if err != nil {
						return "", err
					}
					_, _ = s.(*discordgo.Session).ChannelFileSend(m.(*discordgo.MessageCreate).ChannelID, "me.jpg", f)
				}
			}
		}
		return "", nil
	}
}

func getRandomImage() (string, error) {
	m := make(map[int]string)
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := path.Dir(ex)
	files, err := filepath.Glob(exPath + string(os.PathSeparator) + "*.jpg")
	if err != nil {
		return "", err
	}
	for k, v := range files {
		m[k] = v
	}
	rand.Seed(time.Now().UTC().UnixNano())
	//fmt.Println(rand.Intn(len(m)), len(m))
	str, found := m[rand.Intn(len(m))]
	if !found {
		return "", errors.New("No files found")
	}
	return str, nil
}
