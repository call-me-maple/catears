package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type PatchAlert struct {
	Patch  PatchType
	IDs    *DiscordTrigger
	Offset int
	Args   *PatchOptions
}

func NewPatchAlert() *PatchAlert {
	return &PatchAlert{Args: new(PatchOptions), IDs: new(DiscordTrigger)}
}

type PatchOptions struct {
	Patch     string `arg:"positional"`
	Contract  bool   `arg:"-j,--jane"`
	Stage     uint   `arg:"-s" default:"1"`
	Remainder uint   `arg:"-r" default:"0"`
}

func (pa *PatchAlert) Name() string {
	return pa.Patch.Name()
}

func (pa *PatchAlert) Keywords() []string {
	return AllPatchNames()
}

func (pa *PatchAlert) getIDs() *DiscordTrigger {
	return pa.IDs
}

func (pa *PatchAlert) StatusKey() string {
	keys := []string{pa.IDs.UserID, pa.Patch.Name(), "task"}
	if pa.Args.Contract {
		keys = []string{pa.IDs.UserID, "jane", "task"}
	}
	return formatKey(keys...)
}

func (pa *PatchAlert) NotifyMessage() string {
	content := fmt.Sprintf("<@%v> %v are grown!", pa.IDs.UserID, pa.Patch.Name())
	if pa.Args.Contract {
		content += " Contract is ready!"
	}
	if rand.Intn(10) == 2 || rand.Intn(10) == 4 {
		content += " Tell Jane I say ..hiii.."
	}
	return content
}

func (pa *PatchAlert) NotifyPattern() *regexp.Regexp {
	return regexp.MustCompile(`<@(?P<userId>\d+)> (?P<patch>\w+) are grown!\s?(?P<contract>Contract)?.*`)
}

func (pa *PatchAlert) NotifyParse(m *discordgo.Message) (err error) {
	pa.IDs = triggerFromMessage(m)
	groups := parseNotifier(m.Content, pa)
	pa.Args.Contract = false
	for k, v := range groups {
		switch k {
		case "userId":
			pa.IDs.UserID = v
		case "patch":
			pa.Patch, err = FindPatchType(v)
			if err != nil {
				return err
			}
		case "contract":
			pa.Args.Contract = true
		}
	}
	pa.Args.Stage = 1
	pa.Offset, err = getOffset(pa.IDs.UserID)
	if err != nil {
		return err
	}
	return pa.validate()
}

func (pa *PatchAlert) Parse(m *discordgo.Message) (err error) {
	err = parseMessage(m.Content, pa.Args)
	if err != nil {
		return
	}
	pa.IDs = triggerFromMessage(m)
	pa.Patch, err = FindPatchType(pa.Args.Patch)
	if err != nil {
		return
	}
	pa.Offset, err = getOffset(pa.IDs.UserID)
	if err != nil {
		return
	}
	return pa.validate()
}

func (pa *PatchAlert) validate() error {
	switch {
	// Minus 1 from Stages() because the crops are finished growing at Stages() return value.
	// Need to grow for at least one tick.
	case pa.Args.Stage < 1 || pa.Args.Stage > pa.Patch.Stages()-1:
		return errors.Errorf("Growth stage must be between 1-%v for %v.", pa.Patch.Stages()-1, pa.Patch.Name())
	case pa.Offset < 0 || pa.Offset > 30:
		return errors.Errorf("Farming offset %v out of range 0-30.", pa.Offset)
	default:
		return nil
	}
}

func (pa *PatchAlert) Repeat(mr *discordgo.MessageReactionAdd) error {
	pa.IDs = triggerFromReact(mr)
	return pa.Run()
}

func (pa *PatchAlert) Wait(now time.Time) time.Duration {
	if pa.Args.Remainder == 0 {
		pa.Args.Remainder = pa.Patch.Stages() - pa.Args.Stage
	}

	finish := pa.Patch.getTickTime(int64(pa.Offset), int64(pa.Args.Remainder), now)
	wait := finish.Sub(now)
	return wait
}

func (pa *PatchAlert) FollowUp() time.Duration {
	if pa.Args.Contract {
		return 5 * time.Minute
	}

	wait := pa.Patch.TickRate() / 2
	if wait.Minutes() > 20 {
		wait = 20 * time.Minute
	}
	return wait
}

func (pa *PatchAlert) Run() error {
	return publishAlert(pa)
}
