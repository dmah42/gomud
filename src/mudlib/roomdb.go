package mudlib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var rooms = roomDb{}
var startRoomId string

type roomDb struct {
	memory map[string]*room
}

func init() {
	rooms.memory = make(map[string]*room)
}

// LoadRoomDb loads the room database from the given directory and sets the starting room for new players.
func LoadRoomDb(roomDir, startId string) error {
	startRoomId = startId
	wd, _ := os.Getwd()
	log.Printf("Loading rooms from %s/%s\n", wd, roomDir)
	return filepath.Walk(roomDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return rooms.load(path, fi)
	})
}

func (db roomDb) get(id string) (*room, error) {
	if room, ok := db.memory[id]; ok {
		return room, nil
	}
	return nil, fmt.Errorf("Room %q not found", id)
}

func (db *roomDb) load(path string, fi os.FileInfo) error {
	if fi.IsDir() {
		return nil
	}

	b, err := loadBytes(path)
	if err != nil {
		return err
	}

	if len(b) > 0 {
		newRoom := struct {
			Name        string
			Description string
			Exits       map[string]string
		}{}
		err = json.Unmarshal(b, &newRoom)
		if err != nil {
			return err
		}

		id := path[strings.LastIndex(path, "/")+1:]
		db.memory[id] = &room{
			name:        newRoom.Name,
			description: newRoom.Description,
			exits:       newRoom.Exits,
			playerNicks: make([]string, 0),
		}
		log.Printf("Loaded room %q from %q.\n", id, path)
	}
	return nil
}
