package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/thoas/bokchoy"
)

type HerbOptions struct {
	Stage     uint
	Remainder uint
	ChannelID string
	MessageID string
	UserID    string
}

func NewHerb() *HerbOptions {
	return &HerbOptions{}
}

func (o *HerbOptions) parse(m *discordgo.Message) (err error) {
	buf := new(bytes.Buffer)
	// TODO clean message into str
	str := m.Content
	cmd := flag.NewFlagSet("herb", flag.ContinueOnError)
	cmd.SetOutput(buf)
	cmd.UintVar(&o.Stage, "s", 1, "The current growth stage. 1-4")
	cmd.UintVar(&o.Remainder, "r", 0, "How many growth stages left.")
	err = cmd.Parse(splitCommand(str, "herb"))
	if err != nil {
		err = errors.Errorf("%v", buf.String())
	}
	o.ChannelID = m.ChannelID
	o.MessageID = m.ID
	o.UserID = m.Author.ID
	return
}

func (o *HerbOptions) validate() error {
	switch {
	case o.Stage < 1 || o.Stage > 4:
		return errors.Errorf("Growth stage must be between 1-4.")
	default:
		return nil
	}
}

func (o *HerbOptions) run() (err error) {
	taskKey := formatKey(o.UserID, "herb", "task")
	err = cancelTask(taskKey)
	if err != nil {
		log.Println(err)
		return
	}

	content := fmt.Sprintf("<@%v> Herbs are grown!", o.UserID)

	key := formatKey("config", o.UserID, "offset")
	result, _ := client.Get(key).Result()
	offset, err := strconv.Atoi(result)
	if err != nil {
		_, err = publishMessage(&Message{
			ChannelID:   o.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: "No farm tick offset configured. Use '@catears config offset=value'"},
		})
		if err != nil {
			return
		}
		return
	}
	//var ticks uint
	if o.Remainder == 0 {
		o.Remainder = Herb.Stages() - o.Stage
	}
	finish := getTickTime(int64(offset), 20, int64(o.Remainder))
	wait := time.Until(finish)
	log.Printf("herbs done in: %v at: %v\n", wait, finish)

	task, err := publishMessage(
		&Message{
			ChannelID:   o.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: content},
			Reaction:    "🔁",
			FollowUp: &FollowUp{
				ChannelID: o.ChannelID,
				UserID:    o.UserID,
				Type:      "herbs",
				Key:       taskKey,
				Wait:      10 * time.Minute,
			}},
		bokchoy.WithCountdown(wait))
	if err != nil {
		log.Println("marshal errrrr")
		return
	}
	client.Set(taskKey, task.ID, wait)
	log.Println("set", taskKey, "=", task.ID)

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

//func runHerb(m *discordgo.MessageCreate) (err error) {
//	options := &HerbOptions{
//		Stage:     1,
//		ChannelID: m.ChannelID,
//		MessageID: m.ID,
//		UserID:    m.Author.ID,
//	}
//
//	err = parseHerb(m.Content, options)
//	if err != nil {
//		_, err = publishMessage(&Message{
//			ChannelID:   m.ChannelID,
//			MessageSend: &discordgo.MessageSend{Content: fmt.Sprintf("%v", err)},
//		})
//		if err != nil {
//			return
//		}
//		return
//	}
//	err = sendHerb(options)
//	if err != nil {
//		return
//	}
//	return
//}
