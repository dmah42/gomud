package mudlib

import (
	"reflect"
	"testing"
)

func TestRoomExitDirs(t *testing.T) {
	cases := []struct {
		r         room
		wantExits []string
	}{
		{
			r:         room{name: "r", description: "room", exits: map[string]string{"up": "up"}},
			wantExits: []string{"up"},
		},
		{
			r:         room{name: "r2", description: "room2", exits: map[string]string{}},
			wantExits: []string{},
		},
	}

	for _, tt := range cases {
		gotExits := tt.r.exitDirs()
		if !reflect.DeepEqual(gotExits, tt.wantExits) {
			t.Errorf("want %q, got %q", tt.wantExits, gotExits)
		}
	}
}

func TestRoomDescribe(t *testing.T) {
	cases := []struct {
		r               room
		p               player
		wantDescription string
	}{
		{
			r:               room{name: "room", description: "a room", exits: map[string]string{"a": "a", "b": "b"}, playerNicks: []string{}},
			p:               player{},
			wantDescription: setFgBold(colorGreen, "room\n") + "a room\n" + setFg(colorYellow, "Exits: a, b\n"),
		},
		{
			r:               room{name: "room2", description: "another room", exits: map[string]string{}, playerNicks: []string{"a", "b"}},
			p:               player{nickname: "b"},
			wantDescription: setFgBold(colorGreen, "room2\n") + "another room\na is here.\n",
		},
		{
			r:               room{name: "room2", description: "another room", exits: map[string]string{}, playerNicks: []string{"a", "b"}},
			p:               player{nickname: "c"},
			wantDescription: setFgBold(colorGreen, "room2\n") + "another room\na, b are here.\n",
		},
	}

	for _, tt := range cases {
		gotDescription := tt.r.describe(tt.p)
		if gotDescription != tt.wantDescription {
			t.Errorf("want %q, got %q", tt.wantDescription, gotDescription)
		}
	}
}

func TestRoomAddPlayer(t *testing.T) {
	cases := []struct {
		r           room
		p           string
		wantPlayers []string
	}{
		{r: room{playerNicks: []string{}}, p: "bob", wantPlayers: []string{"bob"}},
		{r: room{playerNicks: []string{"alice"}}, p: "bob", wantPlayers: []string{"alice", "bob"}},
	}

	for _, tt := range cases {
		tt.r.addPlayer(tt.p)
		if !reflect.DeepEqual(tt.r.playerNicks, tt.wantPlayers) {
			t.Errorf("want %v, got %v", tt.wantPlayers, tt.r.playerNicks)
		}
	}
}

func TestRoomRemovePlayer(t *testing.T) {
	cases := []struct {
		r           room
		p           string
		wantPlayers []string
		wantError   bool
	}{
		{r: room{playerNicks: []string{"alice"}}, p: "bob", wantPlayers: []string{"alice"}, wantError: true},
		{r: room{playerNicks: []string{"alice", "bob"}}, p: "bob", wantPlayers: []string{"alice"}, wantError: false},
	}

	for _, tt := range cases {
		gotError := tt.r.removePlayer(tt.p)
		if tt.wantError && gotError == nil {
			t.Errorf("want error, got no error\n")
		}
		if !tt.wantError && gotError != nil {
			t.Errorf("want no error, got error %+v\n", gotError)
		}
		if !reflect.DeepEqual(tt.r.playerNicks, tt.wantPlayers) {
			t.Errorf("want %v, got %v", tt.wantPlayers, tt.r.playerNicks)
		}
	}
}
