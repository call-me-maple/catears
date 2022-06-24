package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/thoas/bokchoy"
	"log"
	"strings"
)

func messageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	data, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = messageCreate.Publish(context.Background(), string(data))
	if err != nil {
		log.Println(err)
		return
	}
}

func reactionCreated(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	data, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = reactionCreate.Publish(context.Background(), string(data))
	if err != nil {
		log.Println(err)
		return
	}
}

func processReaction(r *bokchoy.Request) (err error) {
	res := fmt.Sprintf("%v", r.Task.Payload)
	m := new(discordgo.MessageReactionAdd)
	err = json.Unmarshal([]byte(res), &m)
	if err != nil {
		log.Println(err)
		return err
	}

	// Query full message info
	me, err := dg.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		fmt.Println("failed to grab message", err)
		return
	}

	switch {
	case isBotMentioned(me.Mentions) && strings.Contains(me.Content, "r?") && allReady(me.Reactions):
		err = countDown(me)
	case isBHNotify(me.Content) && isUserMentioned(me.Mentions, m.UserID) && !hasBotReacted(me.Reactions, "✅"):
		err = sendBH(m.ChannelID, m.MessageID, m.UserID, &BHOptions{Seeds: 10})
	case len(me.Embeds) != 0 && m.Emoji.Name == "✅":
		err = archive(m)
	}
	if err != nil {
		log.Println(err)
	}
	return
}

func processMessage(r *bokchoy.Request) (err error) {
	res := fmt.Sprintf("%v", r.Task.Payload)
	m := new(discordgo.MessageCreate)
	err = m.UnmarshalJSON([]byte(res))
	if err != nil {
		log.Println(err)
		return err
	}

	switch {
	case isCommand(m.Content, "music?") && isBotMentioned(m.Mentions):
		err = findMusic(m)
	case isCommand(m.Content, "bh") && isBotMentioned(m.Mentions):
		err = runBirdHouse(m)
	case isBotMentioned(m.Mentions):
		err = respondToMention(m)
	default:
		err = lookForMemes(m)
	}
	if err != nil {
		log.Println(err)
	}
	return
}

func sendMessage(r *bokchoy.Request) (err error) {
	res := fmt.Sprintf("%v", r.Task.Payload)
	m := new(Message)
	err = json.Unmarshal([]byte(res), &m)
	if err != nil {
		log.Println(err)
		return
	}
	sent, err := dg.ChannelMessageSendComplex(m.ChannelID, m.MessageSend)
	if err != nil {
		return err
	}
	log.Println("sent message:", sent.Content)
	reactFollowUp(sent)
	return
}

func sendReaction(req *bokchoy.Request) (err error) {
	res := fmt.Sprintf("%v", req.Task.Payload)
	r := new(Reaction)
	err = json.Unmarshal([]byte(res), &r)
	if err != nil {
		log.Println(err)
		return
	}
	err = dg.MessageReactionAdd(r.ChannelId, r.MessageID, r.Emoji.ID)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("added reaction", r.Emoji.ID)
	return
}
