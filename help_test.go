package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
)

var botID = "979903606901866537"
var allKeywords []string

func init() {
	testDg := &discordgo.Session{
		State: &discordgo.State{
			Ready: discordgo.Ready{
				User: &discordgo.User{ID: botID},
			},
		},
	}
	dg = testDg

	for _, cmd := range initCommands() {
		if kp, ok := cmd.(KeywordProvider); ok {
			allKeywords = append(allKeywords, kp.Keywords()...)
		}
	}
}

func Test_parseNotifier(t *testing.T) {
	type args struct {
		m string
		n Notifier
	}
	tests := []struct {
		args args
		want map[string]string
	}{
		{args{"<@316007672510021634> herbs are grown!", NewPatchAlert()}, map[string]string{"userId": "316007672510021634", "patch": "herbs"}},
		{args{"<@316007672510021634> herbs are grown! Contract is ready!", NewPatchAlert()}, map[string]string{"contract": "Contract", "userId": "316007672510021634", "patch": "herbs"}},
		{args{"<@316007672510021634> herbs are grown! Contract is ready! Tell Jane I say ..hiii..", NewPatchAlert()}, map[string]string{"contract": "Contract", "userId": "316007672510021634", "patch": "herbs"}},
	}

	for _, tt := range tests {
		t.Run("matching notifys", func(t *testing.T) {
			if got := parseNotifier(tt.args.m, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNotifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseMessage(t *testing.T) {
	type args struct {
		m    string
		args interface{}
	}
	tests := []struct {
		args    args
		wantErr error
	}{
		{args{fmt.Sprintf("<@%v> bh -s 9", botID), NewBH()}, nil},
		{args{fmt.Sprintf("<@%v> herb -h", botID), NewBH()}, UserInputError{}},
	}
	for _, tt := range tests {
		t.Run("parsing messagess", func(t *testing.T) {
			if err := parseCommand(tt.args.m, tt.args.args); !errors.Is(err, tt.wantErr) {
				t.Errorf("parseMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_matchesKeyword(t *testing.T) {
	type args struct {
		content string
		mc      MessageCommand
	}
	tests := []struct {
		args args
		want bool
	}{
		{args{fmt.Sprintf("<@%v> bh", botID), NewBH()}, true},
		{args{fmt.Sprintf("<@%v> wow", botID), NewBH()}, false},

		{args{fmt.Sprintf("<@%v> config", botID), NewConfig()}, true},
		{args{fmt.Sprintf("<@%v> wow", botID), NewConfig()}, false},

		{args{fmt.Sprintf("<@%v> d", botID), NewDrop()}, true},
		{args{fmt.Sprintf("<@%v> wow", botID), NewDrop()}, false},

		{args{fmt.Sprintf("<@%v> herb", botID), NewPatchAlert()}, true},
		{args{fmt.Sprintf("<@%v> wow", botID), NewPatchAlert()}, false},

		{args{fmt.Sprintf("<@%v> r?", botID), NewReadyer()}, true},
		{args{fmt.Sprintf("<@%v> rrr", botID), NewReadyer()}, false},
		{args{fmt.Sprintf("<@%v> rrr 100", botID), NewReadyer()}, false},
		{args{fmt.Sprintf("<@%v> ready?", botID), NewReadyer()}, true},
		{args{fmt.Sprintf("<@%v> rrrrrr? 100", botID), NewReadyer()}, false},

		{args{fmt.Sprintf("<@%v> herb things after the word -s 2", botID), NewPatchAlert()}, true},
		{args{"<@10101010101> bh", NewBH()}, false},
		{args{fmt.Sprintf("<@%v>maple", botID), NewPatchAlert()}, false},
		{args{fmt.Sprintf("<@%v> maple   ", botID), NewPatchAlert()}, true},
		{args{fmt.Sprintf("maple <@%v>", botID), NewPatchAlert()}, false},
	}
	for _, tt := range tests {
		t.Run("keyword matching", func(t *testing.T) {
			if got := matchesKeyword(tt.args.content, tt.args.mc); got != tt.want {
				t.Errorf("matchesKeyword(%v %v) = %v, want %v", tt.args.content, tt.args.mc.Keywords()[:5], got, tt.want)
			}
		})
	}
}

func Test_didYouMean(t *testing.T) {
	type args struct {
		search string
		words  []string
	}
	tests := []struct {
		args    args
		want    string
		wantErr error
	}{
		{args{"", []string{}}, "", &EmptySearchError{}},
		{args{"wowee", []string{"xxxxxx", "qpqpqpqpqpq", ""}}, "", &NoSuggestionError{}},
		{args{"birb", NewBH().Keywords()}, "", &MatchingError{}},
		{args{"bir", NewBH().Keywords()}, "Did you mean? bird", nil},
		{args{"d", NewDrop().Keywords()}, "", &MatchingError{}},
		{args{"", allKeywords}, "", &EmptySearchError{}},
	}
	for _, tt := range tests {
		t.Run("mae", func(t *testing.T) {
			got, err := didYouMean(tt.args.search, tt.args.words)
			if err != tt.wantErr {
				t.Errorf("didYouMean() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("didYouMean() = %v, want %v", got, tt.want)
			}
		})
	}
}
