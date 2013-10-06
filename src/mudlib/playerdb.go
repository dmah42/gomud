package mudlib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var players = playerDb{}

type playerDb struct {
	fileMutex sync.Mutex
	filename  string

	memory map[string]player
}

// LoadPlayerDb loads the persistent player database from the given file.
func LoadPlayerDb(filename string) error {
	players.fileMutex.Lock()
	f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		players.fileMutex.Unlock()
		return err
	}
	playerDbLen, err := f.Seek(0, 2)
	if err != nil {
		players.fileMutex.Unlock()
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		players.fileMutex.Unlock()
		return err
	}
	b := make([]byte, playerDbLen)

	_, err = f.Read(b)
	players.fileMutex.Unlock()
	players.filename = filename
	if err != nil {
		return err
	}
	if len(b) > 0 {
		err = json.Unmarshal(b, &players.memory)
		if err != nil {
			return err
		}
	}
	log.Printf("Loaded player database (%d bytes) from %q.\n", len(b), filename)
	return nil
}

func (db *playerDb) get(nickname string) (*player, error) {
	if player, ok := db.memory[nickname]; ok {
		return &player, nil
	}
	return nil, fmt.Errorf("Player %q not found")
}

func (db *playerDb) getAll() []string {
	allPlayers := make([]string, len(db.memory))
	i := 0
	for _, p := range db.memory {
		allPlayers[i] = p.Nickname
		i++
	}
	return allPlayers
}

func (db *playerDb) add(nickname, realname string) (*player, error) {
	if strings.TrimSpace(nickname) == "" {
		return nil, fmt.Errorf("Invalid empty username\n")
	}

	if _, ok := db.memory[nickname]; ok {
		return nil, fmt.Errorf("Player %q already exists", nickname)
	}
	newPlayer := player{
		Nickname: nickname,
		Realname: realname,
	}
	db.memory[nickname] = newPlayer
	if err := db.save(); err != nil {
		log.Printf("Warning: Failed to save player db: %+v\n", err)
	}
	return &newPlayer, nil
}

func (db *playerDb) save() error {
	b, err := json.Marshal(db.memory)
	if err != nil {
		return err
	}
	db.fileMutex.Lock()
	f, err := os.OpenFile(db.filename, os.O_WRONLY, os.ModePerm)
	if err != nil {
		db.fileMutex.Unlock()
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		db.fileMutex.Unlock()
		return err
	}
	db.fileMutex.Unlock()
	log.Printf("Saved player database (%d bytes) to %q.\n", len(b), db.filename)
	return nil
}
