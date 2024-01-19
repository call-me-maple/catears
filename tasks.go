package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	err := cancelTask(alert.StatusKey())
	if err != nil {
		return err
	}
	log.Printf("creating %v alert for user:%v in %v", alert.Name(), alert.getIDs().UserID, alert.Wait(time.Now()))
	task, err := publishMessage(&Message{
		ChannelID:   alert.getIDs().ChannelID,
		MessageSend: &discordgo.MessageSend{Content: alert.NotifyMessage()},
		Reaction:    "🔁",
		FollowUp: &FollowUp{
			ChannelID: alert.getIDs().ChannelID,
			UserID:    alert.getIDs().UserID,
			Name:      alert.Name(),
			Key:       alert.StatusKey(),
			Wait:      alert.FollowUp(),
		}},
		bokchoy.WithCountdown(alert.Wait(time.Now())))
	if err != nil {
		return err
	}

	client.Set(alert.StatusKey(), task.ID, alert.Wait(time.Now()))
	log.Println("set", alert.StatusKey(), "=", task.ID)

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
