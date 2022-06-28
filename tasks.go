package main

import (
	"context"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v7"
	"github.com/thoas/bokchoy"
	"log"
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

func publishMessage(channelID, content string, options ...bokchoy.Option) (task *bokchoy.Task, err error) {
	out := &Message{
		MessageSend: &discordgo.MessageSend{
			Content: content,
		},
		ChannelID: channelID,
	}
	data, err := json.Marshal(out)
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

func publishFollowUp(channelID, userID, taskType, key string, options ...bokchoy.Option) (task *bokchoy.Task, err error) {
	out := &FollowUp{
		ChannelID: channelID,
		UserID:    userID,
		Type:      taskType,
		Key:       key,
	}
	data, err := json.Marshal(out)
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
