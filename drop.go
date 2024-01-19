package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type DropOptions struct {
	IDs  *DiscordTrigger
	Args *DropArgs
}

type DropArgs struct {
	Trigger string `arg:"positional"` // Word thats triggering this cmd
	Length  uint   `arg:"-l" default:"5"`
}

func NewDrop() *DropOptions {
	return &DropOptions{IDs: new(DiscordTrigger), Args: new(DropArgs)}
}

func (o *DropOptions) Name() string {
	return o.Keywords()[0]
}

func (o *DropOptions) getIDs() *DiscordTrigger {
	return o.IDs
}

func (o *DropOptions) StatusKey() string {
	return formatKey(o.IDs.UserID, "drop", "task")
}

func (o *DropOptions) Keywords() []string {
	return []string{"drops", "drop", "d", "meds"}
}

func (o *DropOptions) NotifyMessage() string {
	return fmt.Sprintf("<@%v> %v Placeholder %v :p", o.IDs.UserID, o.Args.Length, o.Name())
}

func (o *DropOptions) NotifyPattern() *regexp.Regexp {
	str := fmt.Sprintf(`<@(?P<userId>\d+)> (?P<length>\d+) Placeholder %v :p`, o.Name())
	return regexp.MustCompile(str)
}

// TODO: This func was copy pasted no cahnge. move somewhere else
func (o *DropOptions) NotifyParse(m *discordgo.Message) (err error) {
	o.IDs = triggerFromMessage(m)
	groups := parseNotifier(m.Content, o)

	for k, v := range groups {
		if k == "userId" {
			o.IDs.UserID = v
		}
		if k == "length" {
			l, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			o.Args.Length = uint(l)
		}
	}

	return o.validate()
}

func (o *DropOptions) Parse(m *discordgo.Message) (err error) {
	err = parseMessage(m.Content, o.Args)
	if err != nil {
		return
	}
	o.IDs = triggerFromMessage(m)
	return o.validate()
}

func (o *DropOptions) validate() (err error) {
	switch {
	case o.Args.Length < 1:
		return errors.Errorf("length must be over 0 :3")
	default:
		return
	}
}

func (o *DropOptions) Repeat(mr *discordgo.MessageReactionAdd) error {
	o.IDs = triggerFromReact(mr)
	return o.Run()
}

func (o *DropOptions) Wait(_ time.Time) time.Duration {
	return time.Duration(o.Args.Length) * time.Hour
}

func (o *DropOptions) FollowUp() time.Duration {
	return 10 * time.Minute
}

func (o *DropOptions) Run() (err error) {
	return publishAlert(o)
}
