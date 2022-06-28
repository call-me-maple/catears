package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/thoas/bokchoy"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func lookForMemes(m *discordgo.MessageCreate) (err error) {
	var content string
	switch strings.ToLower(m.Content) {
	case "ping":
		content = "Pong!"
	case "pong":
		content = "Ping!"
	case "oatmilk":
		content = "don't 🐡"
	case "73":
		content = "fornite"
	default:
		return
	}

	_, err = publishMessage(m.ChannelID, content)
	if err != nil {
		return
	}
	return
}

func respondToMention(m *discordgo.MessageCreate) error {
	var emoji string
	switch {
	case strings.Contains(m.Content, "r?"):
		emoji = "✅"
	case strings.Contains(m.Content, "gm"):
		emoji = "😺"
	default:
		return nil
	}

	_, err := publishReaction(m.ChannelID, m.ID, emoji)
	if err != nil {
		return err
	}
	return nil
}

func findMusic(m *discordgo.MessageCreate) error {
	channels, err := dg.GuildChannels(m.GuildID)
	if err != nil {
		return err
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
	messages, err := dg.ChannelMessages(music.ID, 100, music.LastMessageID, "", "")
	if err != nil {
		fmt.Println("failed to query messages from music", err)
		return err
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
	_, err = publishMessage(m.ChannelID, randMusic.Content)
	if err != nil {
		return err
	}
	return nil
}

func archive(r *discordgo.MessageReactionAdd) (err error) {
	// Find history and backlog channels
	log.Println("in archive")
	channels, err := dg.GuildChannels(r.GuildID)

	if err != nil {
		log.Println(err)
		return
	}
	backlog, err := findChannel(channels, "music")
	if err != nil {
		log.Println(err)
		return
	}
	history, err := findChannel(channels, "history")
	if err != nil {
		log.Println(err)
		return
	}

	// Only move from backlog
	if r.ChannelID != backlog.ID {
		return
	}

	// Check for link in message
	// Query full message info
	me, err := dg.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Println("failed to grab message", err)
		return
	}
	// Delete old message in backlog
	err = dg.ChannelMessageDelete(backlog.ID, r.MessageID)
	if err != nil {
		log.Println("failed to delete old message", err)
		return
	}
	// Create new message in history
	_, err = dg.ChannelMessageSend(history.ID, me.Content)
	if err != nil {
		log.Println("failed to send new message", err)
		return
	}
	return
}

func countDown(m *discordgo.Message) (err error) {
	// Check for a time arg passed in
	re := regexp.MustCompile(`.*r\?.*\s(\d+).*`)
	found := re.FindStringSubmatch(m.Content)
	wait := 3
	if len(found) != 0 {
		intVar, err := strconv.Atoi(found[1])
		if err == nil {
			wait = intVar
		}
	}
	log.Println("counting down", wait)

	for i := wait; i > 0; i-- {
		content := ""
		switch i {
		case 1:
			content = "1 go!"
		case 2, 3, 5, wait:
			content = fmt.Sprintf("%v", i)
		default:
			if i%10 == 0 {
				content = fmt.Sprintf("%v", i)
			}
		}
		if content == "" {
			continue
		}
		delay := bokchoy.WithCountdown(time.Duration(wait-i) * time.Second)
		_, err = publishMessage(m.ChannelID, content, delay)
		if err != nil {
			return
		}
	}

	delay := bokchoy.WithCountdown(time.Duration(rand.Intn((wait/2)+10)) * time.Second)
	_, err = publishMessage(m.ChannelID, "now?", delay)
	if err != nil {
		return
	}

	if rand.Intn(4) != 2 {
		return
	}

	delay = bokchoy.WithCountdown(time.Duration(wait+rand.Intn(15)) * time.Second)
	_, err = publishMessage(m.ChannelID, "now?", delay)
	if err != nil {
		return
	}
	return
}
