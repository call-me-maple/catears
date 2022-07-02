package main

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Message struct {
	ChannelID   string
	MessageSend *discordgo.MessageSend
	FollowUp    *FollowUp
	Reaction    string
}

type FollowUp struct {
	ChannelID string
	UserID    string
	Type      string
	Key       string
	Wait      time.Duration
}
