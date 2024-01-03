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

// TODO publisher interface and one publish function
// or maybe json marshal?
func publishReaction(r *Reaction, options ...bokchoy.Option) (task *bokchoy.Task, err error) {
	data, err := json.Marshal(r)
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

func publishAlert(alert Alerter) error {
	err := cancelTask(alert.getStatusKey())
	if err != nil {
		return err
	}
	log.Printf("creating %v alert for user:%v in %v", alert.getName(), alert.getIDs().UserID, alert.getWait())
	task, err := publishMessage(&Message{
		ChannelID:   alert.getIDs().ChannelID,
		MessageSend: &discordgo.MessageSend{Content: alert.getNotification()},
		Reaction:    "🔁",
		FollowUp: &FollowUp{
			ChannelID: alert.getIDs().ChannelID,
			UserID:    alert.getIDs().UserID,
			Name:      alert.getName(),
			Key:       alert.getStatusKey(),
			Wait:      alert.followUp(),
		}},
		bokchoy.WithCountdown(alert.getWait()))
	if err != nil {
		return err
	}

	client.Set(alert.getStatusKey(), task.ID, alert.getWait())
	log.Println("set", alert.getStatusKey(), "=", task.ID)

	_, err = publishReaction(&Reaction{
		ChannelId: alert.getIDs().ChannelID,
		MessageID: alert.getIDs().MessageID,
		Emoji: &discordgo.Emoji{
			ID: "✅",
		},
	})
	return err
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
