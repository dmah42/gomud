// Binary main defines an mud server.
package main

import (
	"fmt"
	"mudlib"
	"os"
)

const (
	playerDir   = "players/"
	roomDir     = "rooms/"
	configFile  = "config"
	startRoomId = "start"
)

func main() {
	if err := mudlib.LoadPlayerDb(playerDir); err != nil {
		fmt.Println("Failed to load player db: %+v", err)
		os.Exit(1)
	}

	if err := mudlib.LoadRoomDb(roomDir, startRoomId); err != nil {
		fmt.Printf("Failed to load room db: %+v\n", err)
		os.Exit(1)
	}

	if err := mudlib.Run(configFile); err != nil {
		fmt.Printf("Failed to run: %+v\n", err)
		os.Exit(1)
	}
}
