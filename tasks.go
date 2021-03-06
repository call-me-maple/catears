package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v7"
	"github.com/thoas/bokchoy"
)

func cancelTask(key string) error {
	value, err := client.Get(key).Result()
	switch err {
	case nil:
		_, err = messageSend.Cancel(context.Background(), value)
		if err != nil {
			log.Println(err)
		}
	case redis.Nil:
		return nil
	default:
		log.Println("error getting key: ", key, err)
		return err
	}
	return nil
}

//TODO pass Reaction ref instead
func publishReaction(channelID, messageID, emoji string, options ...bokchoy.Option) (task *bokchoy.Task, err error) {
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
	task, err = messageReact.Publish(context.Background(), string(data), options...)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func publishMessage(m *Message, options ...bokchoy.Option) (task *bokchoy.Task, err error) {
	data, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return
	}
	task, err = messageSend.Publish(context.Background(), string(data), options...)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func publishFollowUp(f *FollowUp, options ...bokchoy.Option) (task *bokchoy.Task, err error) {
	data, err := json.Marshal(f)
	if err != nil {
		log.Println(err)
		return
	}
	task, err = followUp.Publish(context.Background(), string(data), options...)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
