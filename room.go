package main

import (
  "fmt"
  "strings"
)

// TODO: add exit directions
type Room struct {
	name  string
	description string
  exitIds []string
  playerNicks []string
}

func (room *Room) ToString() string {
  // TODO: add color. Look up exit names.
  return fmt.Sprintf("%s\n%s\nExits: %s\nPlayers: %s\n", room.name, room.description, strings.Join(room.exitIds, ","), strings.Join(room.playerNicks, ","))
}

func (room *Room) AddPlayer(nick string) {
  room.playerNicks = append(room.playerNicks, nick)
}

func (room *Room) RemovePlayer(nick string) error {
  for i, p := range room.playerNicks {
    if p == nick {
      room.playerNicks = append(room.playerNicks[:i], room.playerNicks[i+1:]...)
      return nil
    }
  }
  return fmt.Errorf("Failed to remove player %q from room %q\n", nick, room.name)
}

