package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/thoas/bokchoy"
)

type ReadyerOptions struct {
	Args *ReadyerArgs
	IDs  *DiscordTrigger
}

type ReadyerArgs struct {
	Trigger string `arg:"positional"` // Word thats triggering this cmd
	Wait    int    `arg:"positional" default:"3"`
	People  int    `arg:"-p" default:"-1"`
}

func NewReadyer() *ReadyerOptions {
	return &ReadyerOptions{IDs: new(DiscordTrigger), Args: new(ReadyerArgs)}
}

func (o *ReadyerOptions) Name() string {
	return o.Keywords()[0]
}

func (o *ReadyerOptions) Keywords() []string {
	return []string{"ready?", "r?", "ready", "r"}
}

func (o *ReadyerOptions) getIDs() *DiscordTrigger {
	return o.IDs
}

func (o *ReadyerOptions) Parse(m *discordgo.Message) (err error) {
	err = parseCommand(m.Content, o.Args)
	if err != nil {
		return
	}
	o.IDs = triggerFromMessage(m)

	channel, err := dg.Channel(o.getIDs().ChannelID)
	if err != nil {
		return
	}
	guild, err := dg.Guild(channel.GuildID)
	if err != nil {
		return
	}
	if o.Args.People < 0 {
		o.Args.People = guild.ApproximatePresenceCount
	}
	return o.validate()
}

func (o *ReadyerOptions) validate() (err error) {
	switch {
	case o.Args.Wait < 0:
		return errors.Errorf("Can't wait negative time.. :p")
	case o.Args.People < 0:
		return errors.Errorf("Can't wait for negative people to ready")
	default:
		return
	}
}

func (o *ReadyerOptions) Run() (err error) {
	m, err := dg.ChannelMessage(o.getIDs().ChannelID, o.getIDs().MessageID)
	if err != nil {
		return
	}

	if !hasBotReacted(m.Reactions, "✅") {
		_, err = publishReaction(&Reaction{
			ChannelId: m.ChannelID,
			MessageID: m.ID,
			Emoji:     &discordgo.Emoji{ID: "✅"}})
		return err
	}

	if r := countReactions(m.Reactions, "✅"); o.Args.People <= r {
		log.Printf("only %v people ready of %v. Nothing to do yet\n", r, o.Args.People)
		return
	}

	log.Println("counting down", o.Args.Wait)
	for i := o.Args.Wait; i > 0; i-- {
		content := ""
		switch i {
		case 1:
			content = "1 go!"
		case 2, 3, 5, o.Args.Wait:
			content = fmt.Sprintf("%v", i)
		default:
			if i%10 == 0 {
				content = fmt.Sprintf("%v", i)
			}
		}
		if content == "" {
			continue
		}
		_, err = publishMessage(
			&Message{
				ChannelID:   o.IDs.ChannelID,
				MessageSend: &discordgo.MessageSend{Content: content},
			},
			bokchoy.WithCountdown(time.Duration(o.Args.Wait-i)*time.Second))

		if err != nil {
			return
		}
	}
	_, err = publishMessage(
		&Message{
			ChannelID:   o.IDs.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: "now?"},
		},
		bokchoy.WithCountdown(time.Duration(rand.Intn(o.Args.Wait*4/3))*time.Second))
	if err != nil {
		return
	}

	if rand.Intn(4) != 2 {
		return
	}

	_, err = publishMessage(
		&Message{
			ChannelID:   o.IDs.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: "i go'ed"},
		},
		bokchoy.WithCountdown(time.Duration(o.Args.Wait+rand.Intn(15))*time.Second))
	if err != nil {
		return
	}
	return
}
