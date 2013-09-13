package main

import (
	"fmt"
)

// TODO: save this to disk
var players = make(map[string]Player)

type Player struct {
	nickname string
	realname string
}

func (player *Player) String() string {
	return fmt.Sprintf("%q (%q)", player.nickname, player.realname)
}

func NewPlayer(nickname, realname string) error {
	if _, ok := players[nickname]; ok {
		return fmt.Errorf("Player %q already exists")
	}
	players[nickname] = Player{
		nickname: nickname,
		realname: realname,
	}
	return nil
}

func GetPlayer(nickname string) (*Player, error) {
	if player, ok := players[nickname]; ok {
		return &player, nil
	}
	return nil, fmt.Errorf("Player %q not found")
}

// TODO: Only return connected players.
func GetAllPlayers() []string {
  allPlayers := make([]string, len(players))
  i := 0
  for  _, p := range players {
     allPlayers[i] = p.nickname
     i++
  }
  return allPlayers
}
