package main

import (
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Runner interface {
	Run() error
}
type Namer interface {
	Name() string
}
type Command interface {
	Runner
	Namer
}
type Parser interface {
	Parse(*discordgo.Message) error
}
type KeywordProvider interface {
	Keywords() []string
}
type DiscordTriggerer interface {
	getIDs() *DiscordTrigger
}
type MessageCommand interface {
	Command
	Parser
	KeywordProvider
	DiscordTriggerer
}
type NotifyPatterner interface {
	NotifyPattern() *regexp.Regexp
}
type NotifyProvider interface {
	NotifyMessage() string
}
type NotifyParser interface {
	NotifyParse(*discordgo.Message) error
}
type Notifier interface {
	NotifyPatterner
	NotifyProvider
	NotifyParser
}
type ReactCommand interface {
	MessageCommand
	Notifier
}

type Waiter interface {
	Wait(time.Time) time.Duration
}
type FollowUper interface {
	FollowUp() time.Duration
}
type Alerter interface {
	Namer
	Statuser
	Notifier
	DiscordTriggerer
	Waiter
	FollowUper
}
type Statuser interface {
	StatusKey() string
}
type Repeater interface {
	Repeat(*discordgo.MessageReactionAdd) error
}

type Tiggerer interface {
	Respond(interface{})
}
