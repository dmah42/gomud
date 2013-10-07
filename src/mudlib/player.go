package mudlib

import (
	"fmt"
	"log"
	"sort"
)

// Just store the nicknames for the connected players
// TODO: we already know which clients are connected.. can we use that instead?
var connected = make([]string, 0)

type player struct {
	Nickname string
	Realname string
	Room     string
}

func (p *player) finger() string {
	return fmt.Sprintf("%s (%s)", p.Nickname, p.Realname)
}

func (p *player) connect() {
	if c, _ := p.isConnected(); c {
		log.Fatalf("User %q is connecting without disconnecting\n", p.Nickname)
	}
	connected = append(connected, p.Nickname)
}

func (p *player) disconnect() {
	if c, index := p.isConnected(); c {
		connected = append(connected[:index], connected[index+1:]...)
		return
	}
	log.Fatalf("User %q is disconnecting without connecting\n", p.Nickname)
}

func (p *player) isConnected() (bool, int) {
	sort.Strings(connected)
	index := sort.SearchStrings(connected, p.Nickname)
	return index < len(connected) && connected[index] == p.Nickname, index
}

func getConnected() []string {
	sort.Strings(connected)
	return connected
}
