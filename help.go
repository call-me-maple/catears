package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type cachedPattern struct {
	regex *regexp.Regexp
	once  sync.Once
	err   error
}

var CachedPatterns = make(map[string]*cachedPattern)

func matchesKeyword(content string, mc MessageCommand) bool {
	name := reflect.TypeOf(mc).String() // is reflect slow? add another func instead?
	p, ok := CachedPatterns[name]
	if !ok {
		p = &cachedPattern{regex: new(regexp.Regexp)}
		CachedPatterns[name] = p
	}

	p.once.Do(func() {
		keys := mc.Keywords()
		for i, key := range keys {
			keys[i] = regexp.QuoteMeta(key)
		}
		str := fmt.Sprintf(`<@(?P<userId>%v)> (?P<trigger>%v)($|\s.*)`, dg.State.User.ID, strings.Join(keys, "|"))
		p.regex, p.err = regexp.Compile(str)
	})

	if p.err != nil {
		log.Printf("error parsing regex? for %v %v\n", name, p.err)
		return false
	}
	return p.regex.MatchString(content)
}

func matchesNotifcation(m *discordgo.Message, np NotifyPatterner) bool {
	return np.NotifyPattern().MatchString(m.Content) &&
		m.Author.ID == dg.State.User.ID
}

type UserInputError struct{}

func (UserInputError) Error() string {
	return "couldnt parse your message.."
}

type NothingTodoError struct{}

func (NothingTodoError) Error() string {
	return "i couldnt find anything to so with ur message"
}

func parseMessage(m *discordgo.Message, mc MessageCommand) error {
	matches := matchesKeyword(m.Content, mc)
	err := mc.Parse(m)
	if err == nil && matches {
		return nil
	}

	if rc, ok := mc.(ReactCommand); ok {
		err = rc.NotifyParse(m)
		matches = matches || matchesNotifcation(m, rc)
	}
	if !matches {
		return NothingTodoError{}
	}
	return err
}

func splitCommand(content, keyword string) []string {
	str := strings.TrimSpace(content)
	_, after, _ := strings.Cut(str, keyword+" ")
	return strings.Split(after, " ")
}

var userMentionRe = regexp.MustCompile(`<@!?(\d+)>`)

func prepCommand(content string) []string {
	// strips any user mentions from message and splits string
	content = userMentionRe.ReplaceAllLiteralString(content, "")
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

type EmptySearchError struct{}

func (EmptySearchError) Error() string {
	return "empty search"
}

type NoSuggestionError struct{}

func (NoSuggestionError) Error() string {
	return "no suggestion"
}

type MatchingError struct{}

func (MatchingError) Error() string {
	return "omg the word was in here"
}

func didYouMean(search string, words []string) (string, error) {
	if search == "" {
		return "", new(EmptySearchError)
	}

	suggestions := fuzzy.RankFind(search, words)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Distance < suggestions[j].Distance
	})
	switch {
	case len(suggestions) < 1 || suggestions[0].Distance >= 10:
		return "", new(NoSuggestionError)
	case suggestions[0].Target == search:
		return "", new(MatchingError)
	default:
		return fmt.Sprintf("Did you mean? %v", suggestions[0].Target), nil
	}
}

func initCommands() (commands []Command) {
	commands = append(commands, NewDrop())
	commands = append(commands, NewPatchAlert())
	commands = append(commands, NewConfig())
	commands = append(commands, NewBH())
	commands = append(commands, NewReadyer())
	return commands
}

func parseCommand(content string, args interface{}) error {
	arg.Parse(args)
	p, err := arg.NewParser(arg.Config{
		IgnoreEnv: true,
		Program:   "catears",
	}, args)
	if err != nil {
		return err
	}
	err = p.Parse(prepCommand(content))
	return err
}

func parseNotifier(content string, n Notifier) map[string]string {
	groups := make(map[string]string)
	r := n.NotifyPattern()
	matches := r.FindStringSubmatch(content)
	if len(r.SubexpNames()) != len(matches) {
		return groups
	}
	for i, name := range r.SubexpNames() {
		// don't add the full match or sub names when they are empty
		if len(name) > 0 && len(matches[i]) > 0 {
			groups[name] = matches[i]
		}
	}
	return groups
}
