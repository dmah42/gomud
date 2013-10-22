package mudlib

import (
	"fmt"
	"strings"
)

type room struct {
	name        string
	description string
	exits       map[string]string
	playerNicks []string
}

func (r room) exitDirs() []string {
	dirs := []string{}
	for k, _ := range r.exits {
		dirs = append(dirs, k)
	}
	return dirs
}

func (r room) describe(p player) string {
	str := setFgBold(colorGreen, fmt.Sprintf("%s\n", r.name))
	str += fmt.Sprintf("%s\n", r.description)
	if len(r.exitDirs()) != 0 {
		exits := strings.Join(r.exitDirs(), ", ")
		str += setFg(colorYellow, fmt.Sprintf("Exits: %s\n", exits))
	}
	playerList := r.playerNicks
	removeStringFromList(p.nickname, &playerList)
	if len(playerList) == 1 {
		str = str + fmt.Sprintf("%s is here.\n", playerList[0])
	} else if len(playerList) > 1 {
		str = str + fmt.Sprintf("%s are here.\n", strings.Join(playerList, ", "))
	}
	return str
}

func (r *room) addPlayer(p string) {
	r.playerNicks = append(r.playerNicks, p)
}

func (r *room) removePlayer(p string) error {
	return removeStringFromList(p, &r.playerNicks)
}
