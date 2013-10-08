package mudlib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var players = playerDb{}

type playerDb struct {
	dir    string
	memory map[string]*player
}

func init() {
	players.memory = make(map[string]*player)
}

// LoadPlayerDb loads the persistent player database from the given directory.
func LoadPlayerDb(playerDir string) error {
	players.dir = playerDir
	wd, _ := os.Getwd()
	log.Printf("Loading players from %s/%s\n", wd, players.dir)
	return filepath.Walk(playerDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return players.load(path, fi)
	})
}

func (db *playerDb) load(path string, fi os.FileInfo) error {
	if fi.IsDir() {
		return nil
	}

	newPlayer := new(player)
	if err := newPlayer.load(path); err != nil {
		return err
	}
	db.memory[newPlayer.nickname] = newPlayer

	return nil
}

func (db playerDb) get(nickname string) (*player, error) {
	if player, ok := db.memory[nickname]; ok {
		return player, nil
	}
	return nil, fmt.Errorf("Player %q not found", nickname)
}

func (db playerDb) getAll() []string {
	allPlayers := make([]string, len(db.memory))
	i := 0
	for _, p := range db.memory {
		allPlayers[i] = p.nickname
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
	newPlayer := &player{
		nickname: nickname,
		realname: realname,
		room:     startRoomId,
	}
	db.memory[nickname] = newPlayer
	newPlayer.save()
	return newPlayer, nil
}
