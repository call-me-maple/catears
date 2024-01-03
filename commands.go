package main

import (
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Runner interface {
	run() error
}
type Namer interface {
	getName() string
}
type Command interface {
	Runner
	Namer
}
type Parser interface {
	parse(*discordgo.Message) error
}
type KeywordProvider interface {
	getKeywords() []string
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
type NotificationPatterner interface {
	getPattern() *regexp.Regexp
}
type NotificationProvider interface {
	getNotification() string
}
type NotificationParser interface {
	parseNotification(*discordgo.Message) error
}
type Notifier interface {
	NotificationPatterner
	NotificationProvider
	NotificationParser
}
type ReactCommand interface {
	MessageCommand
	Notifier
}
type Statuser interface {
	getStatusKey() string
}
type Waiter interface {
	getWait() time.Duration
}
type FollowUper interface {
	followUp() time.Duration
}
type Alerter interface {
	Namer
	Statuser
	Notifier
	DiscordTriggerer
	Waiter
	FollowUper
}
type Repeater interface {
	repeat(*discordgo.MessageReactionAdd) error
}
