package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", os.Getenv("BOT_TOKEN"), "Bot Token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New(fmt.Sprintf("Bot %v", strings.TrimSpace(Token)))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReact)

	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// Move from backlog to history and r? count down
func messageReact(s *discordgo.Session, m *discordgo.MessageReactionAdd) {

	// Only run on check mark emoji
	if m.Emoji.Name != "✅" {
		return
	}

	// Query full message info
	message, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		fmt.Println("failed to grab message", err)
		return
	}

	// Check for 3 ✅ reactions
	var allReady bool
	for _, reaction := range message.Reactions {
		if reaction.Emoji.Name == "✅" {
			// TODO how to test 2 or 3
			allReady = reaction.Count == 3
		}
	}

	for _, mention := range message.Mentions {
		// If all ready, the message mentions the bot, and contains r? then countdown
		if allReady && mention.ID == s.State.User.ID && strings.Contains(message.Content, "r?") {
			// Check for a time arg passed in
			re := regexp.MustCompile(`.*r\?.*\s(\d+).*`)
			found := re.FindStringSubmatch(message.Content)
			waitTime := 3
			if len(found) != 0 {
				intVar, err := strconv.Atoi(found[1])
				if err == nil {
					waitTime = intVar
				}
			}
			go countDown(s, m.ChannelID, waitTime)
		}
	}

	// Find history and backlog channels
	channels, err := s.GuildChannels(m.GuildID)
	if err != nil {
		return
	}
	var backlog *discordgo.Channel
	var history *discordgo.Channel
	for _, channel := range channels {
		if channel.Name == "music" {
			backlog = channel
		}
		if channel.Name == "history" {
			history = channel
		}
	}

	// Only move from backlog
	if m.ChannelID != backlog.ID {
		return
	}

	// Check for link in message
	if len(message.Embeds) == 0 {
		return
	}
	// Delete old message in backlog
	err = s.ChannelMessageDelete(backlog.ID, m.MessageID)
	if err != nil {
		fmt.Println("failed to delete old message", err)
		return
	}
	// Create new message in history
	_, err = s.ChannelMessageSend(history.ID, message.Content)
	if err != nil {
		fmt.Println("failed to send new message", err)
		return
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check that the messages mentions bot
	for _, mention := range m.Mentions {
		if mention.ID != s.State.User.ID {
			continue
		}
		if strings.Contains(m.Content, "r?") {
			s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
		}
		if strings.Contains(m.Content, "gm") {
			s.MessageReactionAdd(m.ChannelID, m.ID, "😺")
		}
		if strings.Contains(m.Content, "music?") {
			go func() {
				channels, err := s.GuildChannels(m.GuildID)
				if err != nil {
					return
				}
				// Find music backlog channels
				var music *discordgo.Channel
				for _, channel := range channels {
					if channel.Name == "music" {
						music = channel
						break
					}
				}
				// Pull messages from music backlog
				messages, err := s.ChannelMessages(music.ID, 100, music.LastMessageID, "", "")
				if err != nil {
					fmt.Println("failed to query messages from music", err)
					return
				}
				// Only recommend messages with links
				var linksOnly []*discordgo.Message
				for _, message := range messages {
					if len(message.Embeds) != 0 {
						linksOnly = append(linksOnly, message)
					}
				}
				// Send a random message with a link
				randMusic := linksOnly[rand.Intn(len(linksOnly))]
				s.ChannelMessageSend(m.ChannelID, randMusic.Content)
			}()
		}
	}

	switch strings.ToLower(m.Content) {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "pong":
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	case "oatmilk":
		s.ChannelMessageSend(m.ChannelID, "don't 🐡")
	case "73":
		s.ChannelMessageSend(m.ChannelID, "fornite")
	}
}

func countDown(s *discordgo.Session, ChannelID string, t int) {
	for i := t; i > 0; i-- {
		switch i {
		case 1:
			s.ChannelMessageSend(ChannelID, "1 go!")
		case 2, 3, 5, t:
			s.ChannelMessageSend(ChannelID, fmt.Sprintf("%v", i))
		default:
			if i%10 == 0 {
				s.ChannelMessageSend(ChannelID, fmt.Sprintf("%v", i))
			}
		}
		time.Sleep(1 * time.Second)
	}
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	s.ChannelMessageSend(ChannelID, "now?")
	time.Sleep(time.Duration(rand.Intn(4)) * time.Second)
	if rand.Intn(6) == 2 {
		s.ChannelMessageSend(ChannelID, "i go'ed")
	}
}
