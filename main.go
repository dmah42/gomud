package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":4242")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msgchan := make(chan Message)
	addchan := make(chan Client)
	rmchan := make(chan Client)

	go handleMessages(msgchan, addchan, rmchan)

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
		conn:     c,
		nickname: promptNick(c, bufc),
		ch:       make(chan Message),
	}
	if strings.TrimSpace(client.nickname) == "" {
		io.WriteString(c, "Invalid Username\n")
		return
	}

	addchan <- client
	defer func() {
		msgchan <- Message{
			nickname:    client.nickname,
			message:     "",
			messageType: messageTypeQuit,
		}
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmchan <- client
	}()
	io.WriteString(c, fmt.Sprintf("Welcome, %s!\n\n", client.nickname))
	msgchan <- Message{
		nickname:    client.nickname,
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

func promptNick(c net.Conn, bufc *bufio.Reader) string {
	io.WriteString(c, addColor(colorBlack, colorRed, "Welcome... to the real world") + "\n")
	io.WriteString(c, "What is your nick? ")
	nick, _, _ := bufc.ReadLine()
	return string(nick)
}
