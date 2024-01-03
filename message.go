package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type FollowUp struct {
	ChannelID string
	UserID    string
	Name      string
	Key       string
	Wait      time.Duration
}

type Message struct {
	ChannelID   string
	MessageSend *discordgo.MessageSend
	FollowUp    *FollowUp
	Reaction    string
}

type DiscordTrigger struct {
	ChannelID string
	MessageID string
	UserID    string
}

func triggerFromMessage(m *discordgo.Message) *DiscordTrigger {
	return &DiscordTrigger{
		ChannelID: m.ChannelID,
		MessageID: m.ID,
		UserID:    m.Author.ID,
	}
}

func triggerFromReact(mr *discordgo.MessageReactionAdd) *DiscordTrigger {
	return &DiscordTrigger{
		ChannelID: mr.ChannelID,
		MessageID: mr.MessageID,
		UserID:    mr.UserID,
	}
}
