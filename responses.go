package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v7"
	"github.com/thoas/bokchoy"
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
	case strings.Contains(m.Content, "r?"):
		out.Emoji.ID = "✅"
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

func countDown(m *discordgo.Message) (err error) {
	// Check for a time arg passed in
	re := regexp.MustCompile(`.*r\?.*\s(\d+).*`)
	found := re.FindStringSubmatch(m.Content)
	wait := 3
	if len(found) != 0 {
		intVar, err := strconv.Atoi(found[1])
		if err == nil {
			wait = intVar
		}
	}
	log.Println("counting down", wait)

	for i := wait; i > 0; i-- {
		content := ""
		switch i {
		case 1:
			content = "1 go!"
		case 2, 3, 5, wait:
			content = fmt.Sprintf("%v", i)
		default:
			if i%10 == 0 {
				content = fmt.Sprintf("%v", i)
			}
		}
		if content == "" {
			continue
		}
		_, err = publishMessage(
			&Message{
				ChannelID:   m.ChannelID,
				MessageSend: &discordgo.MessageSend{Content: content},
			},
			bokchoy.WithCountdown(time.Duration(wait-i)*time.Second))

		if err != nil {
			return
		}
	}
	_, err = publishMessage(
		&Message{
			ChannelID:   m.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: "now?"},
		},
		bokchoy.WithCountdown(time.Duration(rand.Intn(wait*4/3))*time.Second))
	if err != nil {
		return
	}

	if rand.Intn(4) != 2 {
		return
	}

	_, err = publishMessage(
		&Message{
			ChannelID:   m.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: "i go'ed"},
		},
		bokchoy.WithCountdown(time.Duration(wait+rand.Intn(15))*time.Second))
	if err != nil {
		return
	}
	return
}

func sendStatus(m *discordgo.Message, key string) (err error) {
	var content string
	taskID, err := client.Get(key).Result()
	log.Println(key, taskID)
	switch err {
	case nil:
		if task, err := messageSend.Get(context.Background(), taskID); err == nil {
			content = fmt.Sprintf("happening in.. umm... %v", task.ETADisplay())
		} else {
			return err
		}
	case redis.Nil:
		content = "nothing much happening here..."
	default:
		log.Println("error getting key: ", key, err)
		return err
	}
	_, err = publishMessage(
		&Message{
			ChannelID: m.ChannelID,
			MessageSend: &discordgo.MessageSend{
				Content: content,
				Reference: &discordgo.MessageReference{
					MessageID: m.ID,
					ChannelID: m.ChannelID,
				},
			},
		},
	)
	return
}
