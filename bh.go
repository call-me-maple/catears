package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/thoas/bokchoy"
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
		err = publishMessage(m.ChannelID, fmt.Sprintf("%v", err))
		if err != nil {
			return
		}
	}
	err = sendBH(options)
	if err != nil {
		return
	}
	return
}

func sendBH(o *BHOptions) (err error) {
	content := fmt.Sprintf("<@%v> Bird houses are ready!", o.UserID)
	wait := time.Duration(o.Seeds) * 5 * time.Minute
	err = publishMessage(o.ChannelID, content, bokchoy.WithCountdown(wait))
	if err != nil {
		return
	}

	err = publishReaction(o.ChannelID, o.MessageID, "✅")
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
	return
}
