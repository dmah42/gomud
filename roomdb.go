package main

import (
  "encoding/json"
  "fmt"
  "log"
  "os"
  "path/filepath"
  "strings"
  "sync"
)

var roomDb = RoomDb{}

type RoomDb struct {
  fileMutex sync.Mutex
  memory map[string]Room
}

func (db *RoomDb) Get(id string) (*Room, error) {
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

  roomDb.fileMutex.Lock()
  f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
  if err != nil {
    roomDb.fileMutex.Unlock()
    return err
  }
  roomLen, err := f.Seek(0, 2)
  if err != nil {
    roomDb.fileMutex.Unlock()
    return err
  }
  _, err = f.Seek(0, 0)
  if err != nil {
    roomDb.fileMutex.Unlock()
    return err
  }
  b := make([]byte, roomLen)

  _, err = f.Read(b)
  roomDb.fileMutex.Unlock()
  if err != nil {
    return err
  }

  type jsonRoom struct {
    Name string
    Description string
    ExitIds []string
  }

  if len(b) > 0 {
    newRoom := jsonRoom{}
    err = json.Unmarshal(b, &newRoom)
    if err != nil {
      return err
    }
    id := path[strings.LastIndex(path, "/")+1:]
    roomDb.memory[id] = Room {
      name: newRoom.Name,
      description: newRoom.Description,
      exitIds: newRoom.ExitIds,
    }
    log.Printf("Loaded room %q from %q.\n", id, path)
  }
  return nil
}

func LoadRoomDb(roomDir string) error {
  roomDb.memory = make(map[string]Room)
  return filepath.Walk(roomDir, addRoom)
}

