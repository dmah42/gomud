package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type MessageType int

const ( // message types
	messageTypeSay MessageType = iota
	messageTypeEmote
	// TODO: merge join and quit?
	messageTypeJoin
	messageTypeQuit

	messageTypeWho
)

type Message struct {
	nickname    string
	message     string
	messageType MessageType
}

type Client struct {
	conn   net.Conn
	player *Player
	ch     chan Message
}

func (c Client) ReadLinesInto(ch chan<- Message) {
	bufc := bufio.NewReader(c.conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		// TODO: register commands indexed by /<prefix> that create the message to send.
		log.Printf("%q gave command %q.\n", c.player.nickname, line)
		switch {
		// QUIT
		case line == "/quit":
			io.WriteString(c.conn, "Bye!\n")
			c.conn.Close()
      ch <- Message{
        nickname:    c.player.nickname,
        message:     "",
        messageType: messageTypeQuit,
      }
		// EMOTE
		case strings.HasPrefix(line, "/me "):
			ch <- Message{
				nickname:    c.player.nickname,
				message:     line[4:],
				messageType: messageTypeEmote,
			}
			// WHO
		case line == "/who":
			io.WriteString(c.conn, addColor(colorWhite, colorBlack, fmt.Sprintf("%v\n", GetAllPlayers())))
			// FINGER
		case strings.HasPrefix(line, "/finger "):
			player, err := GetPlayer(line[8:])
			if err != nil {
				io.WriteString(c.conn, fmt.Sprintf("%q.\n", err))
			} else {
				io.WriteString(c.conn, addColor(colorWhite, colorBlack, fmt.Sprintf("%v.\n", player)))
			}
		default:
			// SAY
			ch <- Message{
				nickname:    c.player.nickname,
				message:     line,
				messageType: messageTypeSay,
			}
		}
	}
}

func (c Client) WriteLinesFrom(ch <-chan Message) {
	for msg := range ch {
		toPrint := ""
		// TODO: Register command per message type for colors/format string.
		switch {
		case msg.messageType == messageTypeSay:
			toPrint = addColor(colorYellow, colorBlack, fmt.Sprintf("%s says %s", msg.nickname, msg.message))
		case msg.messageType == messageTypeEmote:
			toPrint = addColor(colorGreen, colorBlack, fmt.Sprintf("%s %s", msg.nickname, msg.message))
		case msg.messageType == messageTypeQuit:
			toPrint = addColor(colorRed, colorBlack, fmt.Sprintf("%s has quit.", msg.nickname))
		case msg.messageType == messageTypeJoin:
			toPrint = addColor(colorRed, colorBlack, fmt.Sprintf("%s has joined.", msg.nickname))
		case msg.messageType == messageTypeWho:
		default:
			log.Printf("Unknown message type: %+v", msg)
			return
		}
		_, err := io.WriteString(c.conn, toPrint+"\n")
		if err != nil {
			return
		}
	}
}
