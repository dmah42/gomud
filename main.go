package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

var port = flag.Int("port", 4242, "port to listen on")

const playerDbFilename = "player.db"
const roomDbDir = "rooms/"

func main() {
  err := LoadPlayerDb(playerDbFilename)
  if err != nil {
    fmt.Println("Failed to load player db: %+v", err)
    os.Exit(1)
  }

  err = LoadRoomDb(roomDbDir)
  if err != nil {
    fmt.Printf("Failed to load room db: %+v\n", err)
    os.Exit(1)
  }

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msgchan := make(chan Message)
	addchan := make(chan Client)
	rmchan := make(chan Client)

	go handleMessages(msgchan, addchan, rmchan)

	log.Printf("Listening on " + strconv.Itoa(*port))

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn, msgchan, addchan, rmchan)
	}
}

func handleConnection(c net.Conn, msgchan chan<- Message, addchan chan<- Client, rmchan chan<- Client) {
	bufc := bufio.NewReader(c)
	defer c.Close()

	client := Client{
		conn:   c,
		player: promptNick(c, bufc),
		ch:     make(chan Message),
	}
	if client.player == nil {
		return
	}
	addchan <- client
	msgchan <- Message{
		nickname:    client.player.Nickname,
		message:     "",
		messageType: messageTypeJoin,
	}

  Connect(client.player.Nickname)

	defer func() {
		msgchan <- Message{
		  nickname:    client.player.Nickname,
		  message:     "",
		  messageType: messageTypeQuit,
		}
    Disconnect(client.player.Nickname)
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmchan <- client
	}()

	// TODO: read lines into parser and write response based on location.
	go client.ReadLinesInto(msgchan)
	client.WriteLinesFrom(client.ch)
}

func handleMessages(msgchan <-chan Message, addchan <-chan Client, rmchan <-chan Client) {
	clients := make(map[net.Conn]chan<- Message)
	for {
		select {
		case msg := <-msgchan:
			log.Printf("New message: %+v", msg)
			for _, ch := range clients {
				go func(mch chan<- Message) { mch <- msg }(ch)
			}
		case client := <-addchan:
			log.Printf("New client: %v\n", client.conn)
			clients[client.conn] = client.ch
		case client := <-rmchan:
			log.Printf("Client disconnect: %v\n", client.conn)
			delete(clients, client.conn)
		}
	}
}

func promptNick(c net.Conn, bufc *bufio.Reader) *Player {
	io.WriteString(c, addColor(colorMagenta, colorBlack, "Welcome... to the real world")+"\n")
	errorCount := 0
	var nick, realname string
	for errorCount < 3 {
		io.WriteString(c, "What is your nick? ")
		nickBytes, _, _ := bufc.ReadLine()
		nick = string(nickBytes)
    // TODO: password
    // TODO: check if user is already logged in
    // Check for existing player.
    player, err := playerDb.Get(nick)
    if err == nil {
	    io.WriteString(c, addColor(colorGreen, colorBlack, fmt.Sprintf("Welcome back, %s!\n", nick)))
      log.Printf("Player %+v logged in.\n", player)
      return player
    }
    // Not found so create a new one.
		log.Printf("Creating new player: %s\n", nick)
		io.WriteString(c, "You must be new here. What is your real name? ")
		realnameBytes, _, _ := bufc.ReadLine()
		realname = string(realnameBytes)
		log.Printf("Adding real name %s for %s\n", realname, nick)
		if player, err = playerDb.Add(nick, realname); err == nil {
	    io.WriteString(c, addColor(colorGreen, colorBlack, fmt.Sprintf("Welcome, %s!\n", nick)))
			return player
    } else {
			log.Printf("Error creating new player %s %s: %+v\n", nick, realname, err)
			errorCount = errorCount + 1
		}
	}
	return nil
}
