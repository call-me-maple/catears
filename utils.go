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

func findMentionedUser(mentions []*discordgo.User) (*discordgo.User, error) {
	switch len(mentions) {
	case 1:
		return mentions[0], nil
	case 0:
		return nil, fmt.Errorf("no user mentioned")
	}
	return nil, fmt.Errorf("multiple users mentioned")
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
		if r.Me && r.Emoji.Name == emoji {
			return true
		}
	}
	return false
}
