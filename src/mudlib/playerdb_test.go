package mudlib

import (
	"os"
	"reflect"
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

func TestPlayerDbGet(t *testing.T) {
	if err := LoadPlayerDb("testdata/players"); err != nil {
		t.Fatalf("%+v", err)
	}
	player, err := players.get("alice")
	if err != nil {
		t.Errorf("Failed to get alice")
	}
	if player.realname != "Alice" {
		t.Errorf("want \"Alice\", got %q\n", player.realname)
	}

	player, err = players.get("non-existant")
	if err == nil {
		t.Errorf("Expected error, got %+v", player)
	}
}

func TestPlayerDbGetAll(t *testing.T) {
	want := []string{ "alice", "bob" }
	if err := LoadPlayerDb("testdata/players"); err != nil {
		t.Fatalf("%+v", err)
	}
	got := players.getAll()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %+v, got %+v\n", want, got)
	}
}

func TestPlayerDbAdd(t *testing.T) {
	if err := LoadPlayerDb("testdata/players"); err != nil {
		t.Fatalf("%+v", err)
	}

	player, err := players.get("non-existant")
	if err == nil {
		t.Errorf("Expected error, got %+v", player)
	}

	player, err = players.add("non-existant", "Non Existant")
	if err != nil {
		t.Errorf("%+v", err)
	}

	if player.nickname != "non-existant" || player.realname != "Non Existant" {
		t.Errorf("want %+v, got %+v", []string{"non-existant", "Non Existant"}, player)
	}

	err = os.Remove(players.dir + "/" + player.nickname)
	if err != nil {
		t.Fatalf("Failed to remove new player. Future tests will fail.")
	}
}
