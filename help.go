package main

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/thoas/bokchoy"
	"log"
	"strings"
)

func allReady(users []*discordgo.MessageReactions) bool {
	// TODO look at guild and count active users
	return countReactions(users, "✅") == 3
}

func isCommand(m *discordgo.Message, keyword string) bool {
	keyword = strings.TrimSpace(keyword)
	return (strings.Contains(m.Content, keyword+" ") || strings.HasSuffix(m.Content, keyword)) && isBotMentioned(m.Mentions)
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

func reactFollowUp(m *discordgo.Message) {
	switch {
	case isBHNotify(m, ""):
		err := publishReaction(m.ChannelID, m.ID, "🔁")
		if err != nil {
			return
		}
	case isHerbNotify(m, ""):
		err := publishReaction(m.ChannelID, m.ID, "🔁")
		if err != nil {
			return
		}
	}
}

func publishReaction(channelID, messageID, emoji string, options ...bokchoy.Option) (err error) {
	out := &Reaction{
		ChannelId: channelID,
		MessageID: messageID,
		Emoji: &discordgo.Emoji{
			ID: emoji,
		},
	}

	data, err := json.Marshal(out)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = messageReact.Publish(context.Background(), string(data), options...)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func publishMessage(channelID, content string, options ...bokchoy.Option) (err error) {
	out := &Message{
		MessageSend: &discordgo.MessageSend{
			Content: content,
		},
		ChannelID: channelID,
	}
	data, err := json.Marshal(out)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = messageSend.Publish(context.Background(), string(data), options...)
	if err != nil {
		log.Println(err)
		return err
	}
	return
}
