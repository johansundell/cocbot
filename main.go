package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/johansundell/cocapi"
)

var (
	Token string
	BotID string
)

type commandFunc struct {
	command   string
	helpText  string
	extracter string
}

var botFuncs map[commandFunc]func(string) (string, error) = make(map[commandFunc]func(string) (string, error))
var lockMap = sync.RWMutex{}

var db *sql.DB
var mysqlUser, mysqlPass, mysqlDb, mysqlHost string
var myClanTag, myKey string
var cocClient cocapi.Client

func init() {
	mysqlDb = "cocsniffer"
	mysqlHost = os.Getenv("MYSQL_COC_HOST")
	mysqlUser = os.Getenv("MYSQL_USER")
	mysqlPass = os.Getenv("MYSQL_PASS")
	myClanTag = os.Getenv("COC_CLANTAG")
	myKey = os.Getenv("COC_KEY")
	Token = os.Getenv("DICS_TOKEN")
	emailTo = os.Getenv("EMAIL_TO")
	emailFrom = os.Getenv("EMAIL_FROM")
}

func main() {
	cocClient = cocapi.NewClient(myKey)

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	db, err = sql.Open("mysql", mysqlUser+":"+mysqlPass+"@tcp("+mysqlHost+":3306)/"+mysqlDb)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//fmt.Println(botFuncs[keyTopDonators]("!top donators today"))
	//fmt.Println(getDuration())
	//return

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
		return
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	dg.AddHandler(func(sess *discordgo.Session, event *discordgo.GuildCreate) {
		if event.Guild.Unavailable {
			return
		}

		for _, channel := range event.Guild.Channels {
			if channel.ID == event.Guild.ID {
				channels = append(channels, channel.ID)
				return
			}
		}
	})

	/*dg.AddHandler(func(s *discordgo.Session, event *discordgo.MessageReactionAdd) {
		fmt.Println(event.Emoji)
	})*/

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	go reporter(dg)

	go test(dg)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Author.ID, m.Author.Username, m.Author.Username+"#"+m.Author.Discriminator, m.Content)

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	command := strings.ToLower(m.Content)
	lockMap.RLock()
	defer lockMap.RUnlock()
	msg := ""
	for _, v := range botFuncs {
		if str, err := v(command); err != nil {
			log.Println(err)
		} else {
			msg += str
		}
	}

	// Ugly
	if m.Content == "!send me nude pics" {
		f, err := os.Open("IMG_20170405_142440.jpg")
		if err != nil {
			log.Println(err)
			return
		}
		_, _ = s.ChannelFileSend(m.ChannelID, "me.jpg", f)
	}

	if command == "!help" {
		msg = "**COCBOT COMMANDS**\n```"
		for k, _ := range botFuncs {
			if k.helpText != "" {
				msg += fmt.Sprintf("%s - %s\n", k.command, k.helpText)
			}
		}
		msg += "```"
	}

	if command == "!hidden" && m.Author.Username+"#"+m.Author.Discriminator == "sudde#1958" {
		//s.MessageReactionAdd(m.ChannelID, m.ID, ":raising_hand:")
		msg = "**COCBOT COMMANDS**\n```"
		for k, _ := range botFuncs {
			if k.helpText == "" {
				msg += fmt.Sprintf("%s - %s\n", k.command, k.helpText)
			}
		}
		msg += "```"
	}

	if msg != "" {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}
