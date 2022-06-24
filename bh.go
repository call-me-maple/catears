package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/thoas/bokchoy"
	"strings"
	"time"
)

type BHOptions struct {
	Seeds int
}

func runBirdHouse(m *discordgo.MessageCreate) (err error) {
	options, err := parseBH(m.Content)
	if err != nil {
		err = publishMessage(m.ChannelID, fmt.Sprintf("%v", err))
		if err != nil {
			return
		}
	}
	err = sendBH(m.ChannelID, m.ID, m.Author.ID, options)
	if err != nil {
		return
	}
	return
}

func sendBH(channelID, messageID, userID string, options *BHOptions) (err error) {
	content := fmt.Sprintf("<@%v> Bird houses are ready!", userID)
	wait := time.Duration(options.Seeds) * 5 * time.Minute
	err = publishMessage(channelID, content, bokchoy.WithCountdown(wait))
	if err != nil {
		return
	}

	err = publishReaction(channelID, messageID, "✅")
	if err != nil {
		return
	}
	return
}

func parseBH(str string) (options *BHOptions, err error) {
	options = new(BHOptions)

	str = strings.TrimSpace(str)
	_, after, _ := strings.Cut(str, "bh ")
	input := strings.Split(after, " ")

	buf := new(bytes.Buffer)

	bhCmd := flag.NewFlagSet("bh", flag.ContinueOnError)
	bhCmd.SetOutput(buf)
	bhCmd.IntVar(&options.Seeds, "s", 10, "Number of seeds left.")
	err = bhCmd.Parse(input)
	if err != nil {
		err = errors.Errorf("%v", buf.String())
		return
	}
	return
}
