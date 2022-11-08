package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pborman/getopt/v2"
	"github.com/pkg/errors"
	"github.com/thoas/bokchoy"
)

type ContractOptions struct {
	Stage     uint
	Remainder uint
	Type      PatchType
	ChannelID string
	MessageID string
	UserID    string
}

func runContract(m *discordgo.MessageCreate) (err error) {
	options := &ContractOptions{
		Stage:     1,
		Remainder: 0,
		Type:      Undefined,
		ChannelID: m.ChannelID,
		MessageID: m.ID,
		UserID:    m.Author.ID,
	}

	err = parseContract(m.Content, options)
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

	log.Printf("remainder: %v currentStage: %v\n", options.Remainder, options.Stage)

	err = sendContract(options)
	if err != nil {
		return
	}
	return
}

func sendContract(o *ContractOptions) (err error) {
	taskKey := formatKey(o.UserID, "contract", "task")
	err = cancelTask(taskKey)
	if err != nil {
		log.Println(err)
		return
	}

	content := fmt.Sprintf("<@%v> Contract is ready! Goodluck!", o.UserID)
	if rand.Intn(10) != 2 {
		content += " Tell Jane I say ..hiii.."
	}
	// TODO move to helper getUserOffset
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
	if o.Remainder == 0 {
		o.Remainder = o.Type.Stages() - o.Stage
	}
	log.Println("offset", int64(offset), "tickRate", int64(o.Type.TickRate().Minutes()), "ticks", int64(o.Remainder))
	finish := getTickTime(int64(offset), int64(o.Type.TickRate().Minutes()), int64(o.Remainder))
	wait := time.Until(finish)

	log.Printf("contract done in: %v at: %v\n", wait, finish)

	task, err := publishMessage(
		&Message{
			ChannelID:   o.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: content},
			FollowUp: &FollowUp{
				ChannelID: o.ChannelID,
				UserID:    o.UserID,
				Type:      "contracts",
				Key:       taskKey,
				Wait:      10 * time.Minute,
			}},
		bokchoy.WithCountdown(wait))
	if err != nil {
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

func parseContract(str string, options *ContractOptions) (err error) {
	var opts = getopt.New()
	var s = opts.Uint('s', 1, "The current growth stage.", "stage")
	var r = opts.Uint('r', 0, "How many growth stages left.", "tick")

	str = strings.ToLower(strings.TrimSpace(str))
	i := strings.Index(str, "jane ")
	split := strings.Split(str[i:], " ")

	err = opts.Getopt(split, nil)
	if err != nil {
		buf := new(bytes.Buffer)
		opts.PrintUsage(buf)
		return errors.Errorf("%v\n%v", err, buf.String())
	}
	if opts.NArgs() > 0 {
		options.Type = FindPatchType(opts.Arg(0))
		err = opts.Getopt(opts.Args(), nil)
		if err != nil {
			buf := new(bytes.Buffer)
			opts.PrintUsage(buf)
			return errors.Errorf("%v\n%v", err, buf.String())
		}
	}
	options.Stage = *s
	options.Remainder = *r

	err = options.validate()
	if err != nil {
		return errors.Wrapf(err, "%v", opts.UsageLine())
	}
	return
}

func (o *ContractOptions) validate() (err error) {
	switch {
	case o.Type == Undefined:
		return errors.Errorf("Not sure what your contracts is..")
	case o.Stage < 1 || o.Stage > o.Type.Stages():
		return errors.Errorf("Growth stage must be between 1-%v.", o.Type.Stages())
	}
	return
}
