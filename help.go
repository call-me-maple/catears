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

func isCommand(content, keyword string) bool {
	keyword = strings.TrimSpace(keyword)
	return strings.Contains(content, keyword+" ") || strings.HasSuffix(content, keyword)
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

func isBHNotify(content string) bool {
	return strings.Contains(content, "Bird houses are ready!")
}

func isHerbNotify(content string) bool {
	return strings.Contains(content, "Herbs are grown!")
}

func reactFollowUp(m *discordgo.Message) {
	switch {
	case isBHNotify(m.Content):
		err := publishReaction(m.ChannelID, m.ID, "🔁")
		if err != nil {
			return
		}
	case isHerbNotify(m.Content):
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
