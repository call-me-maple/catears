package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func isBotMentioned(mentions []*discordgo.User) bool {
	return isUserMentioned(mentions, dg.State.User.ID)
}

func isUserMentioned(mentions []*discordgo.User, userID string) bool {
	for _, mention := range mentions {
		if mention.ID == userID {
			return true
		}
	}
	return false
}

func findChannel(channels []*discordgo.Channel, name string) (ch *discordgo.Channel, err error) {
	for _, channel := range channels {
		if channel.Name == name {
			return channel, nil
		}
	}
	return nil, fmt.Errorf("no '%v' channel found", name)
}

func countReactions(mr []*discordgo.MessageReactions, emoji string) (c int) {
	for _, r := range mr {
		if r.Emoji.Name == emoji {
			return r.Count
		}
	}
	return
}
func hasBotReacted(mr []*discordgo.MessageReactions, emoji string) bool {
	for _, r := range mr {
		if r.Emoji.Name == emoji && r.Me {
			return true
		}
	}
	return false
}
