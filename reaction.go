package main

import "github.com/bwmarrin/discordgo"

type Reaction struct {
	ChannelId string
	MessageID string
	Emoji     *discordgo.Emoji
}
