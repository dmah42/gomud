package mudlib

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

// Run listens for connections and handles player interaction.
func Run(port int) error {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		return err
	}

	msgchan := make(chan message)
	addchan := make(chan client)
	rmchan := make(chan client)

	go handleMessages(msgchan, addchan, rmchan)

	log.Printf("Listening on " + strconv.Itoa(port))

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn, msgchan, addchan, rmchan)
	}
}

func handleConnection(c net.Conn, msgchan chan<- message, addchan chan<- client, rmchan chan<- client) {
	bufc := bufio.NewReader(c)
	defer c.Close()

	newClient := client{
		conn:   c,
		player: promptNick(c, bufc),
		ch:     make(chan message),
	}
	if newClient.player == nil {
		return
	}
	addchan <- newClient
	msgchan <- message{
		from:        newClient,
		message:     "",
		messageType: messageTypeJoin,
	}

	newClient.player.connect()

	defer func() {
		msgchan <- message{
			from:        newClient,
			message:     "",
			messageType: messageTypeQuit,
		}
		newClient.player.disconnect()
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmchan <- newClient
	}()

	go newClient.readLinesInto(msgchan)
	newClient.writeLinesFrom(newClient.ch)
}

func handleMessages(msgchan <-chan message, addchan <-chan client, rmchan <-chan client) {
	clients := make(map[net.Conn]chan<- message)
	for {
		select {
		case msg := <-msgchan:
			log.Printf("New message: %+v", msg)
			for _, ch := range clients {
				go func(mch chan<- message) { mch <- msg }(ch)
			}
		case c := <-addchan:
			log.Printf("New client: %v\n", c.conn)
			clients[c.conn] = c.ch
		case c := <-rmchan:
			log.Printf("Client disconnect: %v\n", c.conn)
			delete(clients, c.conn)
		}
	}
}

func promptNick(c net.Conn, bufc *bufio.Reader) *player {
	// TODO: custom prompts
	io.WriteString(c, setFgBold(colorMagenta, "Welcome... to the real world")+"\n")
	errorCount := 0
	var nick, realname string
	for errorCount < 3 {
		io.WriteString(c, "What is your nick? ")
		nickBytes, _, _ := bufc.ReadLine()
		nick = string(nickBytes)
		// TODO: password
		// TODO: check if user is already logged in
		// Check for existing player.
		player, err := players.get(nick)
		if err == nil {
			io.WriteString(c, setFgBold(colorGreen, fmt.Sprintf("Welcome back, %s!\n", nick)))
			// TODO: startup commands
			log.Printf("Player %+v logged in.\n", player)
			return player
		}
		// Not found so create a new one.
		log.Printf("Creating new player: %s\n", nick)
		io.WriteString(c, "You must be new here. What is your real name? ")
		realnameBytes, _, _ := bufc.ReadLine()
		realname = string(realnameBytes)
		log.Printf("Adding real name %s for %s\n", realname, nick)
		if player, err = players.add(nick, realname); err == nil {
			io.WriteString(c, setFgBold(colorGreen, fmt.Sprintf("Welcome, %s!\n", nick)))
			player.Room = startRoomId
			return player
		} else {
			log.Printf("Error creating new player %s %s: %+v\n", nick, realname, err)
			errorCount = errorCount + 1
		}
	}
	return nil
}
