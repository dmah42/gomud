package mudlib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var rooms = roomDb{}
var startRoomId string

type roomDb struct {
	fileMutex sync.Mutex
	memory    map[string]room
}

// LoadRoomDb loads the room database from the given directory and sets the starting room for new players.
func LoadRoomDb(roomDir, startId string) error {
	startRoomId = startId
	rooms.memory = make(map[string]room)
	return filepath.Walk(roomDir, addRoom)
}

func (db *roomDb) get(id string) (*room, error) {
	if room, ok := db.memory[id]; ok {
		return &room, nil
	}
	return nil, fmt.Errorf("Room %q not found", id)
}

func addRoom(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	rooms.fileMutex.Lock()
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		rooms.fileMutex.Unlock()
		return err
	}
	roomLen, err := f.Seek(0, 2)
	if err != nil {
		rooms.fileMutex.Unlock()
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		rooms.fileMutex.Unlock()
		return err
	}
	b := make([]byte, roomLen)

	_, err = f.Read(b)
	rooms.fileMutex.Unlock()
	if err != nil {
		return err
	}

	type jsonRoom struct {
		Name        string
		Description string
		ExitIds     []string
	}

	if len(b) > 0 {
		newRoom := jsonRoom{}
		err = json.Unmarshal(b, &newRoom)
		if err != nil {
			return err
		}
		id := path[strings.LastIndex(path, "/")+1:]
		rooms.memory[id] = room{
			name:        newRoom.Name,
			description: newRoom.Description,
			exitIds:     newRoom.ExitIds,
      playerNicks:   make([]string, 0),
		}
		log.Printf("Loaded room %q from %q.\n", id, path)
	}
	return nil
}
