package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type ConfigOptions struct {
	Key   string
	Value string
	IDs   *DiscordTrigger
}

func NewConfig() *ConfigOptions {
	return &ConfigOptions{IDs: new(DiscordTrigger)}
}

func (o *ConfigOptions) Name() string {
	return o.Keywords()[0]
}

func (o *ConfigOptions) getIDs() *DiscordTrigger {
	return o.IDs
}

func (o *ConfigOptions) Keywords() []string {
	return []string{"config"}
}

func (o *ConfigOptions) Run() (err error) {
	key := strings.Join([]string{"config", o.IDs.UserID, o.Key}, ":")
	client.Set(key, o.Value, 0)
	log.Println("set", key, "=", o.Value)

	reaction := &Reaction{
		ChannelId: o.IDs.ChannelID,
		MessageID: o.IDs.MessageID,
		Emoji: &discordgo.Emoji{
			ID: "✅",
		},
	}
	_, err = publishReaction(reaction)
	return
}

// @catears config key=value
func (options *ConfigOptions) Parse(m *discordgo.Message) (err error) {
	options.IDs.ChannelID = m.ChannelID
	options.IDs.MessageID = m.ID
	options.IDs.UserID = m.Author.ID

	parts := splitCommand(m.Content, "config")
	if len(parts) != 1 {
		return errors.Errorf("Usage '@catears config key=value'")
	}
	parts = strings.Split(parts[0], "=")
	switch len(parts) {
	case 2:
		if !isConfigKey(parts[0]) {
			return errors.Errorf("Not a config options '%v'.", parts[0])
		}
		options.Key, options.Value = parts[0], parts[1]
	default:
		return errors.Errorf("Usage '@catears config key=value'")
	}
	return
}

func (options *ConfigOptions) validate() error {
	// TODO valididate config params
	return nil
}

func getOffset(userID string) (offset int, err error) {
	key := formatKey("config", userID, "offset")
	result, _ := client.Get(key).Result()
	offset, err = strconv.Atoi(result)
	if err != nil {
		err = errors.Errorf("No farm tick offset configured. Use '@catears config offset=value'")
	}
	return
}
