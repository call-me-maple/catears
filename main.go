package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	if Token == "" {
		fmt.Println("no -t arg found checking BOT_TOKEN env variable")
		Token = os.Getenv("BOT_TOKEN")
		fmt.Println("found:", Token)
	}

	dg, err := discordgo.New("Bot " + Token)
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

// Move from backlog to history
func messageReact(s *discordgo.Session, m *discordgo.MessageReactionAdd) {

	// Only run on check mark emoji
	if m.Emoji.Name != "✅" {
		fmt.Println("wrong reaction, not ✅")
		return
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
		fmt.Println("reaction in wrong channel")
		return
	}

	// Check for link in message
	message, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		fmt.Println("failed to grab message")
		return
	}
	if len(message.Embeds) == 0 {
		fmt.Println("no url in message")
		return
	}
	// Delete old message in backlog
	err = s.ChannelMessageDelete(backlog.ID, m.MessageID)
	if err != nil {
		fmt.Printf("failed to delete old message, %v\n", err)
		return
	}
	// Create new message in history
	_, err = s.ChannelMessageSend(history.ID, message.Content)
	if err != nil {
		fmt.Printf("failed to send new message, %v\n", err)
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
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	// You are cute if you are reading this :p
	if m.Content == "oatmilk" {
		s.ChannelMessageSend(m.ChannelID, "don't 🐡")
	}
}
