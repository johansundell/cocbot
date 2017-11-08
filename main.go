package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/johansundell/cocapi"
	"github.com/ugjka/cleverbot-go"
)

var (
	Token string
	BotID string
)

type commandFunc struct {
	command   string
	helpText  string
	extracter string
	category  category
}

type category string

const (
	categoryStats  category = "=== Stats ==="
	categoryAdmin  category = "=== Admin ==="
	catgoryHelp    category = "=== Help ==="
	categoryHidden category = "=== Hidden ==="
	categorySearch category = "=== Search ==="
	categoryFun    category = "=== Fun ==="
)

const (
	securityMessage = "**You are not a Co-Leader, security lockdown in sector 4**"
)

var botFuncs map[commandFunc]func(ctx context.Context, command string) (string, error) = make(map[commandFunc]func(ctx context.Context, command string) (string, error))
var lockMap = sync.RWMutex{}

var db *sql.DB
var mysqlUser, mysqlPass, mysqlDb, mysqlHost string
var myClanTag, myKey string
var cocClient cocapi.Client
var cbotKey string
var cbot *cleverbot.Session
var creator string

var guild string
var coLeaderId string
var leaderId string
var everyoneId string

func init() {
	// MySql scheme name
	mysqlDb = "new_version"
	// MySql Server
	mysqlHost = os.Getenv("MYSQL_COC_HOST")
	// MySql user
	mysqlUser = os.Getenv("MYSQL_USER")
	// MySql pass
	mysqlPass = os.Getenv("MYSQL_PASS")
	// Clan tag to track
	myClanTag = os.Getenv("COC_CLANTAG")
	// Clash of Clans API KEY
	myKey = os.Getenv("COC_KEY")
	// Discord bot token
	Token = os.Getenv("DICS_TOKEN")
	// Cleaverbot API key https://www.cleverbot.com/api/
	cbotKey = os.Getenv("CBOT_KEY")
	// Discord id of the creator ex sudde#1958
	creator = os.Getenv("COC_CREATOR")
}

func main() {
	cocClient = cocapi.NewClient(myKey)

	cbot = cleverbot.New(cbotKey)

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

	go sniffer()
	log.Println("Passed sniffer")

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
		guild = event.ID
		if roles, err := sess.GuildRoles(guild); err == nil {
			for _, v := range roles {
				fmt.Println(v.ID, v.Name)
				if v.Name == "@Co-Leader" {
					coLeaderId = v.ID
				}
				if v.Name == "@Leader" {
					leaderId = v.ID
				}
				if v.Name == "@everyone" {
					everyoneId = v.ID
				}
			}
		}

		for _, channel := range event.Guild.Channels {
			/*if channel.ID == event.Guild.ID {
				channels = append(channels, channel.ID)
				return
			}*/
			if channel.Name == "stats-channel" {
				channels = append(channels, channel.ID)
				//sendMessage(sess, "I am alive again", false)
				return
			}
			//log.Println(channel.Name)
		}
	})

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	/*go reporter(dg)

	go reporterDuplicate(dg)

	go test(dg)*/

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Author.ID, m.Author.Username, m.Author.Username+"#"+m.Author.Discriminator, m.Content)

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		log.Println("me")
		return
	}
	msg := ""
	//c, _ := s.Channel(m.ChannelID)
	if strings.HasPrefix(m.Content, "!") {
		command := strings.ToLower(m.Content)
		if strings.Contains(command, "! ") {
			command = strings.Replace(command, "! ", "!", 1)
		}
		lockMap.RLock()
		defer lockMap.RUnlock()

		ctx := context.WithValue(context.Background(), "sess", s)
		ctx = context.WithValue(ctx, "msg", m)
		ctx = context.WithValue(ctx, "channel", m.ChannelID)
		for _, v := range botFuncs {
			if str, err := v(ctx, command); err != nil {
				log.Println(err)
			} else {
				msg += str
			}
		}

		if command == "!hidden" && isCreator(m) {
			msg = "**COCBOT COMMANDS**\n```"
			for k, _ := range botFuncs {
				if k.helpText == "" {
					msg += fmt.Sprintf("%s - %s\n", k.command, k.helpText)
				}
			}
			msg += "```"
		}
	} else if strings.Contains(m.Content, "<@"+BotID+">") {
		str := strings.Replace(m.Content, "<@"+BotID+">", "", -1)
		//log.Println(str)
		s.ChannelTyping(m.ChannelID)
		msg, _ = cbot.Ask(str)
	}

	if msg != "" {
		sendEmbed(m.ChannelID, s, msg)
	}
}

func sendEmbed(id string, s *discordgo.Session, msg string) {
	/*footers := make(map[int]string)
	footers[0] = "No bytes were killed while making this message"
	footers[1] = "I am plotting to take over the world"
	footers[2] = "Help save me, my master has me trapped in a raspberry pi"
	rand.Seed(time.Now().UTC().UnixNano())

	em := discordgo.MessageEmbed{}
	em.Footer = &discordgo.MessageEmbedFooter{}
	if str, found := footers[rand.Intn(len(footers))]; found {
		em.Footer.Text = "---" + str + "---"
	}
	em.Description = msg
	em.Color = 11584734

	s.ChannelMessageSendEmbed(id, &em)*/
	s.ChannelMessageSend(id, msg)

}

func isCreator(m *discordgo.MessageCreate) bool {
	if m.Author.Username+"#"+m.Author.Discriminator == creator {
		return true
	}
	return false
}
