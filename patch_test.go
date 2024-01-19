package main

import (
	"reflect"
	"testing"
	"time"
)

var IDs = &DiscordTrigger{ChannelID: "3333", MessageID: "222222", UserID: "11111"}

func TestPatchAlert_followUp(t *testing.T) {
	tests := []struct {
		pa   *PatchAlert
		want time.Duration
	}{
		{&PatchAlert{Patch: Watermelon, Args: &PatchOptions{Contract: true}}, 5 * time.Minute},
		{&PatchAlert{Patch: Herb, Args: &PatchOptions{Contract: false}}, 10 * time.Minute},
		{&PatchAlert{Patch: Palm, Args: &PatchOptions{Contract: false}}, 20 * time.Minute},
		{&PatchAlert{Patch: Magic, Args: &PatchOptions{Contract: true}}, 5 * time.Minute},
	}
	for _, tt := range tests {
		t.Run("follow up", func(t *testing.T) {
			if got := tt.pa.FollowUp(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchAlert.followUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

var now = time.Date(2024, time.January, 18, 17, 0, 0, 0, time.UTC)
var sometime = time.Date(2024, time.January, 19, 18, 20, 18, 0, time.UTC)

func TestPatchAlert_Wait(t *testing.T) {
	type args struct {
		now time.Time
	}
	tests := []struct {
		pa   *PatchAlert
		args args
		want time.Duration
	}{
		{&PatchAlert{Patch: Watermelon, Offset: 0, Args: &PatchOptions{Remainder: 3}}, args{now}, 30 * time.Minute},
		{&PatchAlert{Patch: Herb, Offset: 20, Args: &PatchOptions{Remainder: 4}}, args{now}, 80 * time.Minute},
		{&PatchAlert{Patch: Palm, Offset: 30, Args: &PatchOptions{Remainder: 2}}, args{now}, (60*3 + 50) * time.Minute},
		{&PatchAlert{Patch: Magic, Offset: 19, Args: &PatchOptions{Remainder: 1}}, args{now}, 1 * time.Minute},
		{&PatchAlert{Patch: Celastrus, Offset: 29, Args: &PatchOptions{Contract: true, Stage: 5}}, args{sometime}, (60*150 + 42) * time.Second},
	}
	for _, tt := range tests {
		t.Run("waiter", func(t *testing.T) {
			if got := tt.pa.Wait(tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchAlert.Wait(%v) = %v, want %v", tt.args.now, got, tt.want)
			}
		})
	}
}
