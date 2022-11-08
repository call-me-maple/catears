package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type ConfigOptions struct {
	Key       string
	Value     string
	ChannelID string
	MessageID string
	UserID    string
}

func runConfig(m *discordgo.MessageCreate) (err error) {
	options := &ConfigOptions{
		ChannelID: m.ChannelID,
		MessageID: m.ID,
		UserID:    m.Author.ID,
	}

	err = parseConfig(m.Content, options)
	if err != nil {
		_, err = publishMessage(&Message{
			ChannelID:   m.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: fmt.Sprintf("%v", err)},
		})

		if err != nil {
			return
		}
		return nil
	}
	err = saveConfig(options)
	if err != nil {
		return
	}
	return
}

func saveConfig(o *ConfigOptions) (err error) {
	key := strings.Join([]string{"config", o.UserID, o.Key}, ":")
	client.Set(key, o.Value, 0)
	log.Println("set", key, "=", o.Value)

	reaction := &Reaction{
		ChannelId: o.ChannelID,
		MessageID: o.MessageID,
		Emoji: &discordgo.Emoji{
			ID: "✅",
		},
	}
	_, err = publishReaction(reaction)
	if err != nil {
		return
	}
	return
}

// @catears config key=value
func parseConfig(str string, options *ConfigOptions) (err error) {
	parts := splitCommand(str, "config")
	if len(parts) != 1 {
		return errors.Errorf("Usage '@catears config key=value'")
	}
	parts = strings.Split(parts[0], "=")
	switch len(parts) {
	case 2:
		if !isConfigKey(parts[0]) {
			return errors.Errorf("Not a config options '%v'.", parts[0])
		}
		options.Key, options.Value = parts[0], parts[1]
	default:
		return errors.Errorf("Usage '@catears config key=value'")
	}
	return
}
