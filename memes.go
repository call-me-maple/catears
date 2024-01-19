package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func lookForMemes(m *discordgo.Message) (err error) {
	var content string
	switch strings.ToLower(m.Content) {
	case "ping":
		content = "Pong!"
	case "pong":
		content = "Ping!"
	case "oatmilk":
		content = "don't <:bf:979637608718159912>"
	case "73":
		content = "fornite"
	case "barrier":
		content = "nani"
	case "oic":
		content = "<:oic:977400562381365339>"
	case "gm":
		content = "<a:gm:1027902914682961921>"
	default:
		return
	}

	_, err = publishMessage(&Message{
		ChannelID:   m.ChannelID,
		MessageSend: &discordgo.MessageSend{Content: content},
	})
	if err != nil {
		return
	}
	return
}

func respondToMention(m *discordgo.Message) (bool, error) {
	out := &Reaction{
		ChannelId: m.ChannelID,
		MessageID: m.ID,
		Emoji:     &discordgo.Emoji{},
	}

	switch {
	case strings.Contains(m.Content, "gm"):
		out.Emoji.ID = "gm:1027902914682961921"
	default:
		return false, nil
	}

	_, err := publishReaction(out)
	if err != nil {
		return true, err
	}
	return true, nil
}
