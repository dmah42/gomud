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
		errorLog.Fatalf("User %q is connecting without disconnecting\n", p.nickname)
	}
	connected = append(connected, p.nickname)
}

func (p *player) disconnect() {
	if c, index := p.isConnected(); c {
		connected = append(connected[:index], connected[index+1:]...)
		return
	}
	errorLog.Fatalf("User %q is disconnecting without connecting\n", p.nickname)
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
	if err != nil {
		return err
	}

	newPlayer := jsonPlayer{}
	err = json.Unmarshal(b, &newPlayer)
	if err != nil {
		return err
	}

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
		errorLog.Printf("Warning: Failed to save player db: %+v\n", err)
		return
	}
	p.fileMutex.Lock()
	filename := players.dir + "/" + p.nickname
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		p.fileMutex.Unlock()
		errorLog.Printf("Warning: Failed to save player db: %+v\n", err)
		return
	}
	_, err = f.Write(b)
	if err != nil {
		p.fileMutex.Unlock()
		errorLog.Printf("Warning: Failed to save player db: %+v\n", err)
		return
	}
	p.fileMutex.Unlock()
	log.Printf("Saved player %q (%d bytes) to %q.\n", p.nickname, len(b), filename)
}

func (p *player) toRoom(cl client, room string) error {
	currentRoom, err := rooms.get(p.room)
	if err != nil {
		errorLog.Printf("Player %+v is in limbo\n", p)
		return err
	}
	newRoom, err := rooms.get(room)
	if err != nil {
		errorLog.Printf("Player tried to move to unknown room %q -> %q.\n", p.room, room)
		return fmt.Errorf("Moving to unknown room %q.\n", room)
	}
	if p.room != room {
		if err := currentRoom.removePlayer(p.nickname); err != nil {
			errorLog.Printf("%+v", err)
			return err
		}
		msgchan <- message{
			from:        cl,
			message:     room,
			messageType: messageTypeLeaveRoom,
		}
		p.room = room
		p.save()
	}
	newRoom.addPlayer(p.nickname)
	msgchan <- message{
		from:        cl,
		message:     room,
		messageType: messageTypeEnterRoom,
	}
	return nil
}
