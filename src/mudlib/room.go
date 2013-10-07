package mudlib

import (
	"fmt"
	"strings"
)

// TODO: add exit directions
type room struct {
	name        string
	description string
	exitIds     []string
	playerNicks []string
}

func (r *room) describe() string {
	// TODO: add color. Look up exit names.
	str := fmt.Sprintf("%s\n%s\n", r.name, r.description)
	if len(r.exitIds) != 0 {
		str = str + fmt.Sprintf("Exits: %s\n", strings.Join(r.exitIds, ", "))
	}
	if len(r.playerNicks) != 0 {
		str = str + fmt.Sprintf("%s are here.\n", strings.Join(r.playerNicks, ", "))
	}
	return str
}

// TODO: this should take a client, pull out the player, and send messages
func (r *room) addPlayer(nick string) {
	r.playerNicks = append(r.playerNicks, nick)
	// TODO: send message to players in room
}

// TODO: this should take a client, pull out the player, and send messages
func (r *room) removePlayer(nick string) error {
	for i, p := range r.playerNicks {
		if p == nick {
			r.playerNicks = append(r.playerNicks[:i], r.playerNicks[i+1:]...)
			// TODO: send message to players in room
			return nil
		}
	}
	return fmt.Errorf("Failed to remove player %q from room %q\n", nick, r.name)
}
