package main

import (
	"reflect"
	"testing"
	"time"
)

var IDs = &DiscordTrigger{ChannelID: "3333", MessageID: "222222", UserID: "11111"}
var alerts = []PatchAlert{
	PatchAlert{IDs: IDs, Offset: 0, Patch: Herb, Args: &PatchOptions{
		Patch: "herb", Contract: true, Stage: 1, Remainder: 0}},
	PatchAlert{IDs: IDs, Offset: 0, Patch: Flower, Args: &PatchOptions{
		Patch: "flower", Contract: true, Stage: 1, Remainder: 0}},
	PatchAlert{IDs: IDs, Offset: 0, Patch: Whiteberries, Args: &PatchOptions{
		Patch: "wb", Contract: true, Stage: 1, Remainder: 0}},
	PatchAlert{IDs: IDs, Offset: 0, Patch: Watermelon, Args: &PatchOptions{
		Patch: "melon", Contract: true, Stage: 1, Remainder: 0}},
}

func TestPatchAlert_getWait(t *testing.T) {
	tests := []struct {
		pa   *PatchAlert
		want time.Duration
	}{
		{&PatchAlert{Patch: Watermelon, Offset: 0, Args: &PatchOptions{Remainder: 3}}, 5 * time.Minute},
		{&PatchAlert{Patch: Herb, Offset: 20, Args: &PatchOptions{Remainder: 4}}, 10 * time.Minute},
		{&PatchAlert{Patch: Palm, Offset: 30, Args: &PatchOptions{Remainder: 2}}, 20 * time.Minute},
		{&PatchAlert{Patch: Magic, Offset: 19, Args: &PatchOptions{Remainder: 1}}, 5 * time.Minute},
	}
	for _, tt := range tests {
		t.Run("waiter", func(t *testing.T) {
			if got := tt.pa.getWait(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchAlert.getWait() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			if got := tt.pa.followUp(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchAlert.followUp() = %v, want %v", got, tt.want)
			}
		})
	}
}
