package main

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func allReady(users []*discordgo.MessageReactions) bool {
	// TODO look at guild and count active users
	return countReactions(users, "✅") == 3
}

func isCommand(m *discordgo.Message, keyword string) bool {
	keyword = strings.TrimSpace(keyword)
	return (strings.Contains(m.Content, keyword+" ") || strings.HasSuffix(m.Content, keyword)) && isBotMentioned(m.Mentions)
}

func isNotifyCommand(m *discordgo.Message) bool {
	return isCommand(m, "herb") || isCommand(m, "bh")
}

func splitCommand(content, keyword string) []string {
	str := strings.TrimSpace(content)
	_, after, _ := strings.Cut(str, keyword+" ")
	return strings.Split(after, " ")
}

func isConfigKey(key string) bool {
	// TODO: Read keys from redis?
	var keys = []string{"offset"}
	for _, k := range keys {
		if key == k {
			return true
		}
	}
	return false
}

func isNotify(m *discordgo.Message, userID string) (b bool) {
	return isBHNotify(m, userID) || isHerbNotify(m, userID)
}

func isBHNotify(m *discordgo.Message, userID string) (b bool) {
	if userID != "" {
		b = isUserMentioned(m.Mentions, userID)
	} else {
		b = true
	}
	return b && strings.Contains(m.Content, "Bird houses are ready!") && m.Author.ID == dg.State.User.ID
}

func isHerbNotify(m *discordgo.Message, userID string) (b bool) {
	if userID != "" {
		b = isUserMentioned(m.Mentions, userID)
	} else {
		b = true
	}
	return b && strings.Contains(m.Content, "Herbs are grown!") && m.Author.ID == dg.State.User.ID
}

func formatKey(parts ...string) string {
	return strings.Join(parts, ":")
}

func isDev(guildID, channelID string) bool {
	channels, err := dg.GuildChannels(guildID)
	if err != nil {
		return false
	}
	dev, err := findChannel(channels, "dev")
	if err == nil {
		if os.Getenv("ENV") != "DEV" && dev.ID == channelID {
			return true
		}
	}
	return false
}
