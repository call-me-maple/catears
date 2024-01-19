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

type DiscordMessage struct {
	ChannelID   string
	MessageSend *discordgo.MessageSend
	FollowUp    *FollowUp
	Reaction    string
}

func NewDiscordMessage() DiscordMessage {
	return DiscordMessage{
		ChannelID:   "",
		MessageSend: new(discordgo.MessageSend),
		FollowUp:    new(FollowUp),
		Reaction:    "",
	}
}

// todo withOption ... stuff

func (m DiscordMessage) WithReaction(react string) DiscordMessage {
	m.Reaction = react
	return m
}

func (m DiscordMessage) WithFollowUp(followUp *FollowUp) DiscordMessage {
	m.FollowUp = followUp
	return m
}

func (m DiscordMessage) WithMessageSend(send *discordgo.MessageSend) DiscordMessage {
	m.MessageSend = send
	return m
}

func (m DiscordMessage) WithChannelID(id string) DiscordMessage {
	m.ChannelID = id
	return m
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
