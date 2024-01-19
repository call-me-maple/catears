package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v7"
	"github.com/thoas/bokchoy"
)

func messageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't care about the bot's messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	if isDev(m.GuildID, m.ChannelID) {
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
	// Don't care about the bot's reactions
	if r.UserID == s.State.User.ID {
		return
	}
	if isDev(r.GuildID, r.ChannelID) {
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

	for _, c := range initCommands() {
		if mc, ok := c.(MessageCommand); ok && matchesKeyword(me.Content, mc) {
			err = mc.Parse(me)
		}
		if n, ok := c.(Notifier); ok && matchesNotifcation(me, n) {
			err = n.NotifyParse(me)
		}
		if err != nil {
			log.Println(err)
			return
		}

		if r, ok := c.(Repeater); ok && mr.Emoji.Name == "🔁" && isUserMentioned(me.Mentions, mr.UserID) {
			return r.Repeat(mr)
		}
		if s, ok := c.(Statuser); ok && mr.Emoji.Name == "❓" {
			return sendStatus(me, s.StatusKey())
		}
	}
	return
}

func processMessage(r *bokchoy.Request) (err error) {
	res := fmt.Sprintf("%v", r.Task.Payload)

	// Populate MessageCreate with values
	mc := new(discordgo.MessageCreate)
	err = mc.UnmarshalJSON([]byte(res))
	if err != nil {
		log.Println(err)
		return err
	}

	// Query full message info
	m, err := dg.ChannelMessage(mc.ChannelID, mc.ID)
	if err != nil {
		log.Println("failed to grab message", err)
		return
	}

	keywords := []string{}
	for _, c := range initCommands() {
		mc, ok := c.(MessageCommand)
		if !ok {
			continue
		}

		if !matchesKeyword(m.Content, mc) {
			keywords = append(keywords, mc.Keywords()...)
			continue
		}
		err = mc.Parse(m)
		if err != nil {
			_, err = publishMessage(&Message{
				ChannelID:   m.ChannelID,
				MessageSend: &discordgo.MessageSend{Content: fmt.Sprintf("%v", err)},
			})
			return err
		}
		return mc.Run()
	}
	if isBotMentioned(m.Mentions) {
		responded, err := respondToMention(m)
		if err != nil || responded {
			return err
		}
		log.Printf("%v", prepCommand(m.Content))
		guess, err := didYouMean(prepCommand(m.Content)[0], keywords)
		if err != nil {
			return err
		}
		_, err = publishMessage(&Message{
			ChannelID:   m.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: guess},
		})
		if err != nil {
			return err
		}
	}
	return lookForMemes(m)
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

	if m.Reaction != "" {
		reaction := &Reaction{
			ChannelId: m.ChannelID,
			MessageID: sent.ID,
			Emoji: &discordgo.Emoji{
				ID: m.Reaction,
			},
		}

		_, err = publishReaction(reaction)
		if err != nil {
			return
		}
	}
	if m.FollowUp != nil {
		_, err = publishFollowUp(m.FollowUp, bokchoy.WithCountdown(m.FollowUp.Wait))
		if err != nil {
			return
		}
	}
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

func checkin(req *bokchoy.Request) (err error) {
	res := fmt.Sprintf("%v", req.Task.Payload)
	f := new(FollowUp)
	err = json.Unmarshal([]byte(res), &f)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = client.Get(f.Key).Result()
	switch err {
	case nil:
	case redis.Nil:
		content := fmt.Sprintf("I think your %v are still ready <@%v>...", f.Name, f.UserID)
		_, err = publishMessage(&Message{
			ChannelID:   f.ChannelID,
			MessageSend: &discordgo.MessageSend{Content: content},
		})
		return
	default:
		log.Println("error getting key: ", f.Key, err)
		return
	}
	return
}
