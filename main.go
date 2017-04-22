package main

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/johansundell/cocapi"
)

var (
	Token string
	BotID string
)

var buffer = make([][]byte, 0)

type MemberList struct {
	Members []Member
}

type Member struct {
	Name string
	Tag  string
}

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
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
		return
	}

	db, _ = sql.Open("mysql", mysqlUser+":"+mysqlPass+"@tcp("+mysqlHost+":3306)/"+mysqlDb)
	defer db.Close()

	cocClient = cocapi.NewClient(myKey)

	if strings.HasPrefix("!list player #2P9UYQP0", "!list player") {
		//player, _ := cocClient.GetPlayerInfo("#2P9UYQP0")
		//log.Println(player)

		if found, _ := regexp.MatchString("!top war [0-9]+ players", "!top war 1 players"); found {
			str := "!top war 19 players"[len("!top war "):]
			str = str[:strings.Index(str, " ")]
			fmt.Println(strconv.Atoi(str))
		}
	}

	//loadSound()

	// Store the account ID for later use.
	BotID = u.ID
	//fmt.Println(BotID)

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Register guildCreate as a callback for the guildCreate events.
	dg.AddHandler(guildCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
}

// This function will be called (due to AddHandler above) every time a new
// guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			//_, _ = s.ChannelMessageSend(channel.ID, "cocbot is ready, type !hi to say hi to it")
			return
		}
	}
}

func getMembers(search string) (MemberList, error) {
	search = strings.TrimSpace(search)
	log.Println(search)
	rows, err := db.Query("SELECT tag, name FROM members WHERE active = 1 AND name LIKE ? ORDER BY name", "%"+search+"%")
	if err != nil {
		return MemberList{}, err
	}

	mb := MemberList{}
	for rows.Next() {
		m := Member{}
		rows.Scan(&m.Tag, &m.Name)
		mb.Members = append(mb.Members, m)
	}
	return mb, nil
}

func getHeroLvl(p cocapi.Player) int {
	tot := 0
	for _, v := range p.Heroes {
		tot += v.Level
	}
	return tot
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	/*fmt.Println(m.Author.ID, m.Content, m.Mentions)
	for _, v := range m.Mentions {
		if v.ID == BotID {
			msg := ""
			switch m.Content {
			case "hi":
			default:
				msg = "I do not understand " + m.Content
			}

		}
	}*/

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
	//fmt.Println(m.Content)

	msg := ""
	switch {
	case strings.HasPrefix(strings.ToLower(m.Content), "!list members"):
		log.Println("here")
		mb, err := getMembers(m.Content[len("!list members"):])
		if err != nil {
			log.Println(err)
		}
		for _, v := range mb.Members {
			msg += v.Name + " - " + v.Tag + "\n"
		}
	case strings.ToLower(m.Content) == "!show newbie":
		list, _ := cocClient.GetMembers(myClanTag)
		players := []cocapi.Player{}
		for _, v := range list.Items {
			if v.Role == "member" {
				if p, err := cocClient.GetPlayerInfo(v.Tag); err == nil {
					players = append(players, p)
				}
			}
		}
		sort.Slice(players, func(i, j int) bool { return players[i].Donations > players[j].Donations })
		for _, v := range players {
			msg += fmt.Sprintf("%d %s th%d total hero lvl %d\n", v.Donations, v.Name, v.TownHallLevel, getHeroLvl(v))
		}
	case strings.ToLower(m.Content) == "!top donators":
		list, _ := cocClient.GetMembers(myClanTag)
		sort.Slice(list.Items, func(i, j int) bool { return list.Items[i].Donations > list.Items[j].Donations })
		for n := 0; n < len(list.Items); n++ {
			v := list.Items[n]
			msg += fmt.Sprintf("%d %s\n", v.Donations, v.Name)
			if n > 8 {
				break
			}
		}
	case strings.ToLower(m.Content) == "!top war players":
		list, _ := cocClient.GetMembers(myClanTag)
		players := []cocapi.Player{}
		for _, v := range list.Items {
			p, _ := cocClient.GetPlayerInfo(v.Tag)
			players = append(players, p)
		}
		sort.Slice(players, func(i, j int) bool { return players[i].WarStars > players[j].WarStars })
		for i := 0; i < len(players); i++ {
			p := players[i]
			if i > 5 {
				break
			}
			msg += fmt.Sprintf("%d stars for %s\n", p.WarStars, p.Name)
		}
	case strings.ToLower(m.Content) == "!fortune":
		out, err := exec.Command("fortune").Output()
		if err != nil {
			log.Println(err)
			return
		}
		msg = string(out)
	case strings.Contains(strings.ToLower(m.Content), "bbl"):
		msg = "No don't leave me here alone"
	case strings.ToLower(m.Content) == "!help":
		msg = "!list members [name] - to see current members\n!show newbie - to see our newest members and their donations\n!top donators - to see our best donatots\n!top war players - to see our top war whores\n!top war NN players - to see the top NN players\n!send me nude pics - to see me nude"
		//case m.Content == "!show war stars":

	}
	if found, _ := regexp.MatchString("!top war [0-9]+ players", strings.ToLower(m.Content)); found {
		str := strings.ToLower(m.Content)[len("!top war "):]
		str = str[:strings.Index(str, " ")]
		n, _ := strconv.Atoi(str)
		list, _ := cocClient.GetMembers(myClanTag)
		players := []cocapi.Player{}
		for _, v := range list.Items {
			p, _ := cocClient.GetPlayerInfo(v.Tag)
			players = append(players, p)
		}
		sort.Slice(players, func(i, j int) bool { return players[i].WarStars > players[j].WarStars })
		for i := 0; i < len(players); i++ {
			p := players[i]
			if i > n-1 {
				break
			}
			msg += fmt.Sprintf("%d stars for %s\n", p.WarStars, p.Name)
		}
	}

	if msg != "" {
		_, _ = s.ChannelMessageSend(m.ChannelID, msg)
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "!hi" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "testing")
	}

	if m.Content == "!send me nude pics" {
		f, err := os.Open("IMG_20170405_142440.jpg")
		if err != nil {
			log.Println(err)
		}
		//reader := bufio.NewReader()
		_, _ = s.ChannelFileSend(m.ChannelID, "me.jpg", f)
	}

	if m.Content == "!test" {

		return
		f, err := os.Open("/home/johan/Desktop/test.mp3")
		if err != nil {
			log.Println(err)
		}
		//reader := bufio.NewReader()
		_, _ = s.ChannelFileSend(m.ChannelID, "test.mp3", f)
		/*c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			// Could not find guild.
			return
		}
		vc, err := s.ChannelVoiceJoin(g.ID, m.ChannelID, false, true)
		if err != nil {
			return
		}

		// Sleep for a specified amount of time before playing the sound
		time.Sleep(250 * time.Millisecond)

		// Start speaking.
		_ = vc.Speaking(true)

		// Send the buffer data.
		for _, buff := range buffer {
			vc.OpusSend <- buff
		}

		// Stop speaking
		_ = vc.Speaking(false)

		// Sleep for a specificed amount of time before ending.
		time.Sleep(250 * time.Millisecond)

		// Disconnect from the provided voice channel.
		_ = vc.Disconnect()
		log.Println("here")*/
	}

	// If the message is "pong" reply with "Ping!"
	/*if m.Content == "pong" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ping!")
	}*/
}

func loadSound() error {
	file, err := os.Open("/home/johan/Desktop/test.dca")

	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}
		fmt.Println("test", opuslen)
		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}
