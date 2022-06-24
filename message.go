package main

import (
	"github.com/bwmarrin/discordgo"
)

type Message struct {
	ChannelID   string
	MessageSend *discordgo.MessageSend
}
