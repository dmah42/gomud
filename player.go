package main

import (
	"fmt"
  "log"
  "sort"
)

// Just store the nicknames for the connected players
var connected = make([]string, 0)

type Player struct {
	Nickname string
	Realname string
  Room  string
}

func (player *Player) String() string {
	return fmt.Sprintf("%q (%q)", player.Nickname, player.Realname)
}

func Connect(nickname string) {
  if IsConnected(nickname) {
    log.Fatalf("User %q is connecting without disconnecting\n", nickname)
  }
  connected = append(connected, nickname)
}

func Disconnect(nickname string) {
  sort.Strings(connected)
  index := sort.SearchStrings(connected, nickname)
  if index == len(connected) || connected[index] != nickname {
    log.Fatalf("User %q is disconnecting without connecting\n", nickname)
  }
  connected = append(connected[:index], connected[index + 1:]...)
}

func IsConnected(nickname string) bool {
  sort.Strings(connected)
  index := sort.SearchStrings(connected, nickname)
  return index < len(connected) && connected[index] == nickname
}

func GetConnected() []string {
  sort.Strings(connected)
  return connected
}

