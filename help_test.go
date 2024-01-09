package main

import (
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func Test_isCommand(t *testing.T) {
	type args struct {
		m       *discordgo.Message
		keyword string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"words before", args{&discordgo.Message{Content: "@<979903606901866537> hii ce red"}, "drop"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCommand(tt.args.m, tt.args.keyword); got != tt.want {
				t.Errorf("isCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchesKeyword(t *testing.T) {
	type args struct {
		m  *discordgo.Message
		kp KeywordProvider
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"red?", args{&discordgo.Message{Content: "@<979903606901866537> red"}, &DropOptions{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesKeyword(tt.args.m, tt.args.kp); got != tt.want {
				t.Errorf("matchesKeyword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseNotifier(t *testing.T) {
	type args struct {
		m *discordgo.Message
		n Notifier
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"dont include empty matches", args{&discordgo.Message{Content: "<@316007672510021634> herbs are grown!"}, NewPatchAlert()}, map[string]string{"userId": "316007672510021634", "patch": "herbs"}},
		{"contract parse", args{&discordgo.Message{Content: "<@316007672510021634> herbs are grown! Contract is ready!"}, NewPatchAlert()}, map[string]string{"contract": "Contract", "userId": "316007672510021634", "patch": "herbs"}},
		{"eggy contract parse", args{&discordgo.Message{Content: "<@316007672510021634> herbs are grown! Contract is ready! Tell Jane I say ..hiii.."}, NewPatchAlert()}, map[string]string{"contract": "Contract", "userId": "316007672510021634", "patch": "herbs"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNotifier(tt.args.m, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNotifier() = %v, want %v", got, tt.want)
			}
		})
	}
}
