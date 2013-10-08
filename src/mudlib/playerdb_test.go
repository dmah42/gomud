package mudlib

import (
	"testing"
)

func TestPlayerDbLoad(t *testing.T) {
	if err := LoadPlayerDb("testdata/players"); err != nil {
		t.Fatalf("%+v", err)
	}
	player, ok := players.memory["alice"]
	if !ok {
		t.Errorf("Failed to get alice")
	}
	if player.realname != "Alice" {
		t.Errorf("want \"Alice\", got %q\n", player.realname)
	}
}

func TestPlayerDbAdd(t *testing.T) {
	// TODO
}

func TestPlayerDbGet(t *testing.T) {
	// TODO
}

func TestPlayerDbGetAll(t *testing.T) {
	// TODO
}
