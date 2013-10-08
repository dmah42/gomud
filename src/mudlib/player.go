package mudlib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

// Just store the nicknames for the connected players
// TODO: we already know which clients are connected.. can we use that instead?
var connected = make([]string, 0)

// TODO: save player when updating Room

type jsonPlayer struct {
	Nickname string
	Realname string
	Room     string
}

type player struct {
	nickname  string
	realname  string
	room      string
	fileMutex sync.Mutex
	filename  string
}

func (p *player) finger() string {
	return fmt.Sprintf("%s (%s)", p.nickname, p.realname)
}

func (p *player) connect() {
	if c, _ := p.isConnected(); c {
		log.Fatalf("User %q is connecting without disconnecting\n", p.nickname)
	}
	connected = append(connected, p.nickname)
}

func (p *player) disconnect() {
	if c, index := p.isConnected(); c {
		connected = append(connected[:index], connected[index+1:]...)
		return
	}
	log.Fatalf("User %q is disconnecting without connecting\n", p.nickname)
}

func (p *player) isConnected() (bool, int) {
	sort.Strings(connected)
	index := sort.SearchStrings(connected, p.nickname)
	return index < len(connected) && connected[index] == p.nickname, index
}

func getConnected() []string {
	sort.Strings(connected)
	return connected
}

func (p *player) load(path string) error {
  b, err := loadBytes(path)
  if err != nil { return err }

	newPlayer := jsonPlayer{}
	err = json.Unmarshal(b, &newPlayer)
	if err != nil { return err }

	p.nickname = newPlayer.Nickname
	p.realname = newPlayer.Realname
	p.room = newPlayer.Room
	log.Printf("Loaded player %q from %q.\n", p.nickname, path)
	return nil
}

func (p player) save() {
	jp := jsonPlayer{
		Nickname: p.nickname,
		Realname: p.realname,
		Room:     p.room,
	}
	b, err := json.Marshal(jp)
	if err != nil {
		log.Printf("Warning: Failed to save player db: %+v\n", err)
		return
	}
	p.fileMutex.Lock()
	filename := players.dir + p.nickname
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		p.fileMutex.Unlock()
		log.Printf("Warning: Failed to save player db: %+v\n", err)
		return
	}
	_, err = f.Write(b)
	if err != nil {
		p.fileMutex.Unlock()
		log.Printf("Warning: Failed to save player db: %+v\n", err)
		return
	}
	p.fileMutex.Unlock()
	log.Printf("Saved player %q (%d bytes) to %q.\n", p.nickname, len(b), filename)
}
