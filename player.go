package main

import (
  "encoding/json"
	"fmt"
  "log"
  "os"
  "strings"
)

const playerDb = "player.db"

var players = make(map[string]Player)

type Player struct {
	Nickname string
	Realname string
}

func (player *Player) String() string {
	return fmt.Sprintf("%q (%q)", player.Nickname, player.Realname)
}

func LoadPlayerDb() error {
  f, err := os.OpenFile(playerDb, os.O_RDONLY | os.O_CREATE, os.ModePerm)
  if err != nil {
    return err
  }
  dbLen, err := f.Seek(0, 2)
  if err != nil {
    return err
  }
  _, err = f.Seek(0, 0)
  if err != nil {
    return err
  }
  b := make([]byte, dbLen)
  _, err = f.Read(b)
  if err != nil {
    return err
  }
  if len(b) > 0 {
    err = json.Unmarshal(b, &players)
    if err != nil {
      return err
    }
  }
  log.Printf("Loaded player database (%d bytes) from %q.\n", len(b), playerDb)
  return nil
}

func SavePlayerDb() error {
  b, err := json.Marshal(players)
  if err != nil {
    return err
  }
  f, err := os.OpenFile(playerDb, os.O_WRONLY, os.ModePerm)
  if err != nil {
    return err
  }
  _, err = f.Write(b)
  if err != nil {
    return err
  }
  log.Printf("Saved player database (%d bytes) to %q.\n", len(b), playerDb)
	return nil
}

func NewPlayer(nickname, realname string) error {
	if strings.TrimSpace(nickname) == "" {
		return fmt.Errorf("Invalid empty username\n")
	}

	if _, ok := players[nickname]; ok {
		return fmt.Errorf("Player %q already exists", nickname)
	}
	players[nickname] = Player{
		Nickname: nickname,
		Realname: realname,
	}
  return SavePlayerDb()
}

func GetPlayer(nickname string) (*Player, error) {
	if player, ok := players[nickname]; ok {
		return &player, nil
	}
	return nil, fmt.Errorf("Player %q not found")
}

// TODO: Only return connected players.
func GetAllPlayers() []string {
  allPlayers := make([]string, len(players))
  i := 0
  for  _, p := range players {
     allPlayers[i] = p.Nickname
     i++
  }
  return allPlayers
}
