package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type BHOptions struct {
	Args *BHArgs
	IDs  *DiscordTrigger
}

type BHArgs struct {
	Trigger string `arg:"positional"` // Word thats triggering this cmd
	Seeds   int    `arg:"-s" default:"10"`
}

func NewBH() *BHOptions {
	return &BHOptions{IDs: new(DiscordTrigger), Args: new(BHArgs)}
}

func (o *BHOptions) getName() string {
	return "bird houses"
}

func (o *BHOptions) getKeywords() []string {
	return []string{"bh", "bird", "house", "birdhouse", "birdhouses", "birb"}
}

func (o *BHOptions) getIDs() *DiscordTrigger {
	return o.IDs
}

func (o *BHOptions) getStatusKey() string {
	return formatKey(o.IDs.UserID, "bh", "task")
}

func (o *BHOptions) getNotification() string {
	return fmt.Sprintf("<@%v> %v are ready!", o.IDs.UserID, o.getName())
}

func (o *BHOptions) getPattern() *regexp.Regexp {
	str := fmt.Sprintf(`<@(?P<userId>\d+)> %v are ready!`, o.getName())
	return regexp.MustCompile(str)
}

func (o *BHOptions) parseNotification(m *discordgo.Message) (err error) {
	o.IDs = triggerFromMessage(m)
	groups := parseNotifier(m, o)

	for k, v := range groups {
		if k == "userId" {
			o.IDs.UserID = v
		}
	}

	return o.validate()
}

func (o *BHOptions) parse(m *discordgo.Message) (err error) {
	err = parseMessage(m, o.Args)
	if err != nil {
		return
	}
	o.IDs = triggerFromMessage(m)
	return o.validate()
}

func (options *BHOptions) validate() (err error) {
	switch {
	case options.Args.Seeds < 0 || options.Args.Seeds > 10:
		return errors.Errorf("Seeds must be between 1-10.")
	default:
		return
	}
}

func (o *BHOptions) repeat(mr *discordgo.MessageReactionAdd) error {
	o.IDs = triggerFromReact(mr)
	o.Args.Seeds = 10
	return o.run()
}

func (o *BHOptions) getWait() time.Duration {
	return time.Duration(o.Args.Seeds) * 5 * time.Minute
}

func (o *BHOptions) followUp() time.Duration {
	return 5 * time.Minute
}

func (o *BHOptions) run() (err error) {
	return publishAlert(o)
}
