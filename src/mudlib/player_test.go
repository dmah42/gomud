package mudlib

import (
	"reflect"
	"testing"
)

func TestPlayerFinger(t *testing.T) {
	cases := []struct {
		p    player
		want string
	}{
		{
			p:    player{Nickname: "nick", Realname: "real", Room: "room"},
			want: "nick (real)",
		},
		{
			p:    player{Nickname: "nick_name", Realname: "real name", Room: "room"},
			want: "nick_name (real name)",
		},
	}

	for _, tt := range cases {
		got := tt.p.finger()
		if got != tt.want {
			t.Errorf("want %q, got %q", tt.want, got)
		}
	}
}

func TestPlayerIsConnected(t *testing.T) {
	cases := []struct {
		p           player
		wantConnect bool
		wantIndex   int
	}{
		{
			p:           player{Nickname: "b", Realname: "bob", Room: "beehive"},
			wantConnect: true,
			wantIndex:   1,
		},
		{
			p:           player{Nickname: "a", Realname: "alice", Room: "aviary"},
			wantConnect: true,
			wantIndex:   0,
		},
		{
			p:           player{Nickname: "c", Realname: "charles", Room: "chapel"},
			wantConnect: false,
			wantIndex:   2,
		},
	}
	for _, tt := range cases {
		if tt.wantConnect {
			connected = append(connected, tt.p.Nickname)
		}
	}

	for _, tt := range cases {
		gotConnect, gotIndex := tt.p.isConnected()
		if gotConnect != tt.wantConnect || gotIndex != tt.wantIndex {
			t.Errorf("wantConnect %v, gotConnect %v, wantIndex %v, gotIndex %v", tt.wantConnect, gotConnect, tt.wantIndex, gotIndex)
		}
	}
	connected = make([]string, 0)
}

func TestPlayerConnect(t *testing.T) {
	cases := []struct {
		p           player
		wantConnect bool
	}{
		{
			p:           player{Nickname: "a", Realname: "alice", Room: "aviary"},
			wantConnect: true,
		},
		{
			p:           player{Nickname: "b", Realname: "bob", Room: "beehive"},
			wantConnect: false,
		},
	}
	for _, tt := range cases {
		if tt.wantConnect {
			tt.p.connect()
		}
	}

	for _, tt := range cases {
		gotConnect, _ := tt.p.isConnected()
		if gotConnect != tt.wantConnect {
			t.Errorf("wantConnect %v, gotConnect %v", tt.wantConnect, gotConnect)
		}
	}
	connected = make([]string, 0)
}

func TestPlayerDisconnect(t *testing.T) {
	cases := []struct {
		p              player
		wantDisconnect bool
	}{
		{
			p:              player{Nickname: "a", Realname: "alice", Room: "aviary"},
			wantDisconnect: true,
		},
		{
			p:              player{Nickname: "b", Realname: "bob", Room: "beehive"},
			wantDisconnect: false,
		},
	}
	for _, tt := range cases {
		tt.p.connect()
		if tt.wantDisconnect {
			tt.p.disconnect()
		}
	}

	for _, tt := range cases {
		gotConnect, _ := tt.p.isConnected()
		if gotConnect == tt.wantDisconnect {
			t.Errorf("wantDisconnect %v, gotDisconnect %v", tt.wantDisconnect, !gotConnect)
		}
	}
	connected = make([]string, 0)
}

func TestPlayerGetConnected(t *testing.T) {
	cases := []struct {
		p                []player
		wantGetConnected []string
	}{
		{
			p: []player{
				player{Nickname: "b", Realname: "bob", Room: "beehive"},
				player{Nickname: "a", Realname: "alice", Room: "aviary"},
				player{Nickname: "c", Realname: "charles", Room: "chapel"},
			},
			wantGetConnected: []string{"a", "b", "c"},
		},
	}
	for _, tt := range cases {
		for _, p := range tt.p {
			p.connect()
		}
		gotGetConnected := getConnected()
		if !reflect.DeepEqual(gotGetConnected, tt.wantGetConnected) {
			t.Errorf("wantGetConnected %q, gotGetConnected %q", tt.wantGetConnected, gotGetConnected)
		}
	}
	connected = make([]string, 0)
}
