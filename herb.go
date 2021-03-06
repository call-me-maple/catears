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

func runHerb(m *discordgo.MessageCreate) (err error) {
	options := &HerbOptions{
		Stage:     1,
		ChannelID: m.ChannelID,
		MessageID: m.ID,
		UserID:    m.Author.ID,
	}

	err = parseHerb(m.Content, options)
	if err != nil {
		_, err = publishMessage(&Message{
			ChannelID:   m.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: fmt.Sprintf("%v", err)},
		})
		if err != nil {
			return
		}
		return
	}
	err = sendHerb(options)
	if err != nil {
		return
	}
	return
}

func sendHerb(o *HerbOptions) (err error) {
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

	_, err = publishReaction(o.ChannelID, o.MessageID, "✅")
	if err != nil {
		return
	}
	return
}

func parseHerb(str string, options *HerbOptions) (err error) {
	buf := new(bytes.Buffer)

	cmd := flag.NewFlagSet("herb", flag.ContinueOnError)
	cmd.SetOutput(buf)
	cmd.UintVar(&options.Stage, "s", 1, "The current growth stage. 1-4")
	cmd.UintVar(&options.Remainder, "r", 0, "How many growth stages left.")
	err = cmd.Parse(splitCommand(str, "herb"))
	if err != nil {
		err = errors.Errorf("%v", buf.String())
		return
	}

	return options.validate()
}

func (options *HerbOptions) validate() (err error) {
	switch {
	case options.Stage < 1 || options.Stage > 4:
		return errors.Errorf("Growth stage must be between 1-4.")
	}
	return
}
