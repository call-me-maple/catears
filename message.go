package main

import (
	"github.com/bwmarrin/discordgo"
)

type Message struct {
	ChannelID   string
	MessageSend *discordgo.MessageSend
}

type FollowUp struct {
	ChannelID string
	UserID    string
	Type      string
	Key       string
}
