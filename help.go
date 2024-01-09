package main

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/pkg/errors"
)

func allReady(users []*discordgo.MessageReactions) bool {
	// TODO look at guild and count active users
	if os.Getenv("ENV") == "DEV" {
		return true
	}
	return countReactions(users, "✅") == 3
}

func matchesKeyword(m *discordgo.Message, kp KeywordProvider) bool {
	for _, keyword := range kp.getKeywords() {
		if isCommand(m, keyword) {
			return true
		}
	}
	return false
}

func matchesNotifcation(m *discordgo.Message, np NotificationPatterner) bool {
	return np.getPattern().MatchString(m.Content) &&
		m.Author.ID == dg.State.User.ID
}

func isCommand(m *discordgo.Message, keyword string) bool {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	content := strings.ToLower(m.Content)
	return (strings.Contains(content, keyword+" ") && isBotMentioned(m.Mentions))
}

func splitCommand(content, keyword string) []string {
	str := strings.TrimSpace(content)
	_, after, _ := strings.Cut(str, keyword+" ")
	return strings.Split(after, " ")
}
func prepCommand(m *discordgo.Message) []string {
	// strips any user mentions from message and splits string
	content := m.Content
	for _, user := range m.Mentions {
		content = strings.NewReplacer(
			"<@"+user.ID+">", "",
			"<@!"+user.ID+">", "",
		).Replace(content)
	}
	return strings.Split(strings.TrimSpace(content), " ")
}

func isConfigKey(key string) bool {
	// TODO: Read keys from redis?
	var keys = []string{"offset"}
	for _, k := range keys {
		if key == k {
			return true
		}
	}
	return false
}

func formatKey(parts ...string) string {
	return strings.Join(parts, ":")
}

func isDev(guildID, channelID string) bool {
	dev := os.Getenv("ENV") != "DEV"

	channels, err := dg.GuildChannels(guildID)
	if err != nil {
		return false
	}
	c, err := findChannel(channels, "dev")
	if err != nil {
		log.Println(err)
		return true
	}

	return dev == (c.ID == channelID)
}

func didYouMean(search string, words []string) error {
	suggestions := fuzzy.RankFind(search, words)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Distance < suggestions[j].Distance
	})
	switch {
	case len(suggestions) > 0 && suggestions[0].Distance < 10:
		return errors.Errorf("Did you mean? %v", suggestions[0].Target)
	default:
		return nil
	}
}

func initCommands() (commands []Command) {
	commands = append(commands, NewDrop())
	commands = append(commands, NewPatchAlert())
	commands = append(commands, NewConfig())
	commands = append(commands, NewBH())
	return commands
}

func parseMessage(m *discordgo.Message, args interface{}) error {
	arg.Parse(args)
	p, err := arg.NewParser(arg.Config{
		IgnoreEnv: true,
		Program:   "catears",
	}, args)
	if err != nil {
		return err
	}
	err = p.Parse(prepCommand(m))
	return err
}

func parseNotifier(m *discordgo.Message, n Notifier) map[string]string {
	groups := make(map[string]string)
	r := n.getPattern()
	matches := r.FindStringSubmatch(m.Content)
	for i, name := range r.SubexpNames() {
		// don't add the full match or sub names when they are empty
		if len(name) > 0 && len(matches[i]) > 0 {
			groups[name] = matches[i]
		}
	}
	return groups
}
