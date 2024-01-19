package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type BHOptions struct {
	Trigger string          `arg:"positional"` // Word thats triggering this cmd
	Seeds   int             `arg:"-s" default:"10"`
	IDs     *DiscordTrigger `arg:"-"` // interface trigger with function to respond with output
}

func NewBH() *BHOptions {
	return &BHOptions{IDs: new(DiscordTrigger)}
}

func (o *BHOptions) Name() string {
	return "bird houses"
}

func (o *BHOptions) Keywords() []string {
	return []string{"bh", "bird", "house", "birdhouse", "birdhouses", "birb"}
}

func (o *BHOptions) getIDs() *DiscordTrigger {
	return o.IDs
}

func (o *BHOptions) StatusKey() string {
	return formatKey(o.IDs.UserID, "bh", "task")
}

func (o *BHOptions) NotifyMessage() string {
	return fmt.Sprintf("<@%v> %v are ready!", o.IDs.UserID, o.Name())
}

func (o *BHOptions) NotifyPattern() *regexp.Regexp {
	str := fmt.Sprintf(`<@(?P<userId>\d+)> %v are ready!`, o.Name())
	return regexp.MustCompile(str)
}

func (o *BHOptions) NotifyParse(m *discordgo.Message) (err error) {
	o.IDs = triggerFromMessage(m)
	groups := parseNotifier(m.Content, o)

	for k, v := range groups {
		if k == "userId" {
			o.IDs.UserID = v
		}
	}

	return o.validate()
}

func (o *BHOptions) Parse(m *discordgo.Message) (err error) {
	err = parseCommand(m.Content, o)
	if err != nil {
		return
	}
	o.IDs = triggerFromMessage(m)
	return o.validate()
}

func (o *BHOptions) validate() (err error) {
	switch {
	case o.Seeds < 0 || o.Seeds > 10:
		return errors.Errorf("Seeds must be between 1-10.")
	default:
		return
	}
}

func (o *BHOptions) Repeat(mr *discordgo.MessageReactionAdd) error {
	o.IDs = triggerFromReact(mr)
	o.Seeds = 10
	return o.Run()
}

func (o *BHOptions) Wait(_ time.Time) time.Duration {
	return time.Duration(o.Seeds) * 5 * time.Minute
}

func (o *BHOptions) FollowUp() time.Duration {
	return 5 * time.Minute
}

func (o *BHOptions) Run() (err error) {
	return publishAlert(o)
}
