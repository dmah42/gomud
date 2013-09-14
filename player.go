package main

import (
	"fmt"
)

// TODO: track connected players

type Player struct {
	Nickname string
	Realname string
}

func (player *Player) String() string {
	return fmt.Sprintf("%q (%q)", player.Nickname, player.Realname)
}

