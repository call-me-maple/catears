package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/thoas/bokchoy"
	"log"
	"time"
)

type BHOptions struct {
	Seeds     int
	ChannelID string
	MessageID string
	UserID    string
}

func runBirdHouse(m *discordgo.MessageCreate) (err error) {
	options := &BHOptions{
		Seeds:     10,
		ChannelID: m.ChannelID,
		MessageID: m.ID,
		UserID:    m.Author.ID,
	}

	err = parseBH(m.Content, options)
	if err != nil {
		_, err = publishMessage(m.ChannelID, fmt.Sprintf("%v", err))
		if err != nil {
			return
		}
		return
	}
	err = sendBH(options)
	if err != nil {
		return
	}
	return
}

func sendBH(o *BHOptions) (err error) {
	taskKey := formatKey("bh", o.UserID, "task")
	err = cancelTask(taskKey)
	if err != nil {
		log.Println(err)
		return
	}

	content := fmt.Sprintf("<@%v> Bird houses are ready!", o.UserID)
	wait := time.Duration(o.Seeds) * 5 * time.Minute
	task, err := publishMessage(o.ChannelID, content, bokchoy.WithCountdown(wait))
	if err != nil {
		return
	}
	client.Set(taskKey, task.ID, wait)
	log.Println("set", taskKey, "=", task.ID)

	checkup := wait + (5 * time.Minute)
	_, err = publishFollowUp(o.ChannelID, o.UserID, "bird houses", taskKey, bokchoy.WithCountdown(checkup))
	if err != nil {
		return err
	}

	_, err = publishReaction(o.ChannelID, o.MessageID, "✅")
	if err != nil {
		return
	}
	return
}

func parseBH(str string, options *BHOptions) (err error) {
	buf := new(bytes.Buffer)
	cmd := flag.NewFlagSet("bh", flag.ContinueOnError)
	cmd.SetOutput(buf)
	cmd.IntVar(&options.Seeds, "s", 10, "Number of seeds left.")
	err = cmd.Parse(splitCommand(str, "bh"))
	if err != nil {
		err = errors.Errorf("%v", buf.String())
		return
	}
	return options.validate()
}

func (options *BHOptions) validate() (err error) {
	switch {
	case options.Seeds > 0 && options.Seeds <= 10:
	default:
		return errors.Errorf("Seeds must be between 1-10.")
	}
	return
}
