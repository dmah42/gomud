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
	"strings"
)

var port = flag.Int("port", 4242, "port to listen on")

func main() {
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
	if strings.TrimSpace(client.player.nickname) == "" {
		io.WriteString(c, "Invalid Username\n")
		return
	}

	addchan <- client
	defer func() {
		msgchan <- Message{
			nickname:    client.player.nickname,
			message:     "",
			messageType: messageTypeQuit,
		}
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmchan <- client
	}()
	io.WriteString(c, fmt.Sprintf("Welcome, %s!\n\n", client.player.nickname))
	msgchan <- Message{
		nickname:    client.player.nickname,
		message:     "",
		messageType: messageTypeJoin,
	}

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
	var nick string
	for errorCount < 3 {
		io.WriteString(c, "What is your nick? ")
		nickBytes, _, _ := bufc.ReadLine()
		nick = string(nickBytes)
		log.Printf("Creating new player: %s\n", nick)
		// TODO: prompt for realname
		if err := NewPlayer(nick, "realname"); err != nil {
			// TODO: check password
			log.Printf("Error creating new player %s: %+v\n", nick, err)
			errorCount = errorCount + 1
		} else {
			// TODO: check error
			player, _ := GetPlayer(nick)
			return player
		}
	}
	return nil
}
