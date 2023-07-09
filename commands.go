package main

import (
	"log"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/bwmarrin/discordgo"
)

type Parser interface {
	parse(*discordgo.Message) error
}
type Runner interface {
	run() error
}
type Validator interface {
	validate() error
}
type Command interface {
	Parser
	Runner
	Validator
}

type PatchAlert struct {
	Patch     PatchType
	ChannelID string
	MessageID string
	UserID    string
	Args      *PatchArgs
}

func NewPatch() *PatchAlert {
	return &PatchAlert{}
}

type PatchArgs struct {
	Patch     string `arg:"positional"`
	Contract  bool   `arg:"-j,--jane"`
	Stage     uint   `arg:"-s"`
	Remainder uint   `arg:"-r"`
}

func (pa *PatchAlert) parse(m *discordgo.Message) error {
	var args PatchArgs
	arg.Parse(&args)
	p, err := arg.NewParser(arg.Config{
		IgnoreEnv: true,
		Program:   "catears",
	}, &args)

	if err != nil {
		log.Println(err)
		return err
	}

	// move to help.go
	content := m.Content
	for _, user := range m.Mentions {
		content = strings.NewReplacer(
			"<@"+user.ID+">", "",
			"<@!"+user.ID+">", "",
		).Replace(content)
	}
	argv := strings.Split(strings.TrimSpace(content), " ")

	err = p.Parse(argv)
	if err != nil {
		log.Println(err)
		return err
	}

	patch, err := FindPatchType(args.Patch)
	patch.TickRate()
	if err != nil {
		return err
	}
	pa.ChannelID = m.ChannelID
	pa.MessageID = m.ID
	pa.UserID = m.Author.ID
	pa.Patch = patch
	pa.Args = &args

	return nil
}
