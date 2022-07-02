package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gorhill/cronexpr"
	"github.com/pkg/errors"
	"github.com/thoas/bokchoy"
)

type HerbOptions struct {
	Stage     uint
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
	// All times ending in :00, :20, :40
	timings := [3]int{00, 20, 40}
	for i, t := range timings {
		t -= offset
		if t < 0 {
			t += 60
		}
		timings[i] = t
	}
	expr := fmt.Sprintf("0 %v,%v,%v * ? * * *", timings[0], timings[1], timings[2])
	parse, err := cronexpr.Parse(expr)
	if err != nil {
		return err
	}

	// 5 growth stages for herbs
	growthTimes := parse.NextN(time.Now(), 5-o.Stage)
	finish := growthTimes[len(growthTimes)-1]
	wait := finish.Sub(time.Now())
	log.Printf("herbs done in: %v at: %v\n", wait, growthTimes[len(growthTimes)-1])
	task, err := publishMessage(
		&Message{
			ChannelID:   o.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: content}},
		bokchoy.WithCountdown(wait))
	if err != nil {
		return
	}
	client.Set(taskKey, task.ID, wait)
	log.Println("set", taskKey, "=", task.ID)

	checkup := wait + (10 * time.Minute)
	_, err = publishFollowUp(o.ChannelID, o.UserID, "herbs", taskKey, bokchoy.WithCountdown(checkup))
	if err != nil {
		return err
	}

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
	err = cmd.Parse(splitCommand(str, "herb"))
	if err != nil {
		err = errors.Errorf("%v", buf.String())
		return
	}
	return options.validate()
}

func (options *HerbOptions) validate() (err error) {
	switch {
	case options.Stage >= 1 && options.Stage < 5:
	default:
		return errors.Errorf("Growth stage must be between 1-4.")
	}
	return
}
