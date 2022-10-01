package main

import (
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func allReady(users []*discordgo.MessageReactions) bool {
	// TODO look at guild and count active users
	return countReactions(users, "✅") == 3
}

func isCommand(m *discordgo.Message, keyword string) bool {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	content := strings.ToLower(m.Content)
	return (strings.Contains(content, keyword+" ") || strings.HasSuffix(content, keyword)) && isBotMentioned(m.Mentions)
}

func isNotifyCommand(m *discordgo.Message) bool {
	return isCommand(m, "herb") || isCommand(m, "bh") || isCommand(m, "jane") || isCommand(m, "d1") || isCommand(m, "d4")
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
	return isBHNotify(m, userID) || isHerbNotify(m, userID) || isContractNotify(m, userID) || isDrop1Notify(m, userID) || isDrop4Notify(m, userID)
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

func isDrop1Notify(m *discordgo.Message, userID string) (b bool) {
	if userID != "" {
		b = isUserMentioned(m.Mentions, userID)
	} else {
		b = true
	}
	return b && strings.Contains(m.Content, "Placeholder drop 1!") && m.Author.ID == dg.State.User.ID
}

func isDrop4Notify(m *discordgo.Message, userID string) (b bool) {
	if userID != "" {
		b = isUserMentioned(m.Mentions, userID)
	} else {
		b = true
	}
	return b && strings.Contains(m.Content, "Placeholder drop 4!") && m.Author.ID == dg.State.User.ID
}

func isContractNotify(m *discordgo.Message, userID string) (b bool) {
	if userID != "" {
		b = isUserMentioned(m.Mentions, userID)
	} else {
		b = true
	}
	return b && strings.Contains(m.Content, "Contract is ready! Goodluck!") && m.Author.ID == dg.State.User.ID
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

func getTickTime(offset, tickRate, ticks int64) time.Time {
	calcOffset := (offset % tickRate) * 60
	unixNow := time.Now().Unix() + calcOffset

	currentTick := (unixNow - (unixNow % (tickRate * 60)))
	goalTick := currentTick + (ticks * tickRate * 60)

	return time.Unix(goalTick-calcOffset, 0)
}
