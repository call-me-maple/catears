package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/thoas/bokchoy"
)

type DropOptions struct {
	ChannelID string
	MessageID string
	UserID    string
	Length    int
}

func parseDrop(m *discordgo.MessageCreate) (err error) {
	if len(m.Content) < 1 || strings.ToLower(m.Content)[0] != 'd' {

	}
	return
}

func runDrop1(m *discordgo.MessageCreate) (err error) {
	options := &DropOptions{
		ChannelID: m.ChannelID,
		MessageID: m.ID,
		UserID:    m.Author.ID,
		Length:    3,
	}

	err = sendDrop(options)
	if err != nil {
		return
	}
	return
}

func runDrop4(m *discordgo.MessageCreate) (err error) {
	options := &DropOptions{
		ChannelID: m.ChannelID,
		MessageID: m.ID,
		UserID:    m.Author.ID,
		Length:    4,
	}

	err = sendDrop(options)
	if err != nil {
		return
	}
	return
}

func sendDrop(o *DropOptions) (err error) {
	taskKey := formatKey(o.UserID, fmt.Sprintf("drop%v", o.Length), "task")
	err = cancelTask(taskKey)
	if err != nil {
		log.Println(err)
		return
	}

	content := fmt.Sprintf("<@%v> Placeholder drop %v!", o.UserID, o.Length)
	wait := time.Duration(o.Length) * time.Hour

	task, err := publishMessage(
		&Message{
			ChannelID:   o.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: content},
			Reaction:    "🔁",
			FollowUp: &FollowUp{
				ChannelID: o.ChannelID,
				UserID:    o.UserID,
				Type:      "drop",
				Key:       taskKey,
				Wait:      5 * time.Minute,
			}},
		bokchoy.WithCountdown(wait))
	if err != nil {
		return
	}
	client.Set(taskKey, task.ID, wait)
	log.Println("set", taskKey, "=", task.ID)
	reaction := &Reaction{
		ChannelId: o.ChannelID,
		MessageID: o.MessageID,
		Emoji: &discordgo.Emoji{
			ID: "✅",
		},
	}
	_, err = publishReaction(reaction)
	if err != nil {
		return
	}
	return
}
