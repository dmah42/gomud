package mudlib

import (
	"testing"
)

func TestRoomDescribe(t *testing.T) {
	cases := []struct {
		r    room
		wantDescription string
	}{
		{
			r:    room{name: "room", description: "a room", exitIds: []string{"a", "b"}, playerNicks: []string{}},
      wantDescription: "room\na room\nExits: a, b\n",
		},
		{
			r:    room{name: "room2", description: "another room", exitIds: []string{}, playerNicks: []string{"a", "b"},},
      wantDescription: "room2\nanother room\na, b are here.\n",
		},
	}

	for _, tt := range cases {
		gotDescription := tt.r.describe()
		if gotDescription != tt.wantDescription {
			t.Errorf("want %q, got %q", tt.wantDescription, gotDescription)
		}
	}
}

func TestRoomAddPlayer(t *testing.T) {
// TODO
}

func TestRoomRemovePlayer(t *testing.T) {
// TODO
}
