// Binary main defines an mud server.
package main

import (
	"flag"
	"fmt"
	"mudlib"
	"os"
)

var port = flag.Int("port", 4242, "port to listen on")

const (
	playerDir   = "players/"
	roomDir     = "rooms/"
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

	if err := mudlib.Run(*port); err != nil {
		os.Exit(1)
	}
}
