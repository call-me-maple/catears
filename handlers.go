package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/thoas/bokchoy"
	"log"
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
	mr := new(discordgo.MessageReactionAdd)
	err = json.Unmarshal([]byte(res), &mr)
	if err != nil {
		log.Println(err)
		return err
	}

	// Query full message info
	me, err := dg.ChannelMessage(mr.ChannelID, mr.MessageID)
	if err != nil {
		fmt.Println("failed to grab message", err)
		return
	}

	switch {
	// Count-down ready check
	case isCommand(me, "r?") && allReady(me.Reactions):
		err = countDown(me)
	// Bird House repeat check
	case isBHNotify(me, mr.UserID) && mr.Emoji.Name == "🔁" && !hasBotReacted(me.Reactions, "✅"):
		err = sendBH(&BHOptions{
			Seeds:     10,
			ChannelID: mr.ChannelID,
			MessageID: mr.MessageID,
			UserID:    mr.UserID})
	// Herb repeat check
	case isHerbNotify(me, mr.UserID) && mr.Emoji.Name == "🔁" && !hasBotReacted(me.Reactions, "✅"):
		err = sendHerb(&HerbOptions{
			Stage:     0,
			ChannelID: mr.ChannelID,
			MessageID: mr.MessageID,
			UserID:    mr.UserID,
		})
	// Archive check
	case len(me.Embeds) != 0 && mr.Emoji.Name == "✅":
		err = archive(mr)
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
	// Query full message info
	me, err := dg.ChannelMessage(m.ChannelID, m.ID)
	if err != nil {
		fmt.Println("failed to grab message", err)
		return
	}

	switch {
	case isCommand(me, "music?"):
		err = findMusic(m)
	case isCommand(me, "bh"):
		err = runBirdHouse(m)
	case isCommand(me, "herb"):
		err = runHerb(m)
	case isCommand(me, "config"):
		err = runConfig(m)
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
		log.Printf("%v\n%v\n", err, m)
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
