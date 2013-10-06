// Package mudlib is a mud engine.
package mudlib

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type messageType int

const ( // message types
	messageTypeSay messageType = iota
	messageTypeEmote
	// TODO: merge join and quit?
	messageTypeJoin
	messageTypeQuit

	messageTypeWho
)

type message struct {
	nickname    string
	message     string
	messageType messageType
}

type client struct {
	conn   net.Conn
	player *player
	ch     chan message
}

func (c client) readLinesInto(ch chan<- message) {
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
		log.Printf("%q gave command %q.\n", c.player.Nickname, line)
		switch {
		// QUIT
		case line == "/quit":
			io.WriteString(c.conn, "Bye!\n")
			c.conn.Close()
			ch <- message{
				nickname:    c.player.Nickname,
				message:     "",
				messageType: messageTypeQuit,
			}
		// EMOTE
		case strings.HasPrefix(line, "/me "):
			ch <- message{
				nickname:    c.player.Nickname,
				message:     line[4:],
				messageType: messageTypeEmote,
			}
			// WHO
		case line == "/who":
			io.WriteString(c.conn, addColor(colorWhite, colorBlack, fmt.Sprintf("%v\n", getConnected())))
			// FINGER
		case strings.HasPrefix(line, "/finger "):
			if player, err := players.get(line[8:]); err == nil {
				toPrint := addColor(colorWhite, colorBlack, fmt.Sprintf("%+v ", player.finger()))
				if c,_ := player.isConnected(); c {
					toPrint += addColor(colorGreen, colorBlack, "[online]\n")
				} else {
					toPrint += addColor(colorRed, colorBlack, "[offline]\n")
				}
				io.WriteString(c.conn, toPrint)
			} else {
				io.WriteString(c.conn, fmt.Sprintf("%q.\n", err))
			}
		case line == "look":
			room, err := rooms.get(c.player.Room)
			if err == nil {
				io.WriteString(c.conn, room.String())
			} else {
				// TODO: handle limbo
				io.WriteString(c.conn, fmt.Sprintf("%q.\n", err))
				log.Printf("%q in limbo %q.\n", c.player.Nickname, c.player.Room)
			}
		default:
			// SAY
			ch <- message{
				nickname:    c.player.Nickname,
				message:     line,
				messageType: messageTypeSay,
			}
		}
	}
}

func (c client) writeLinesFrom(ch <-chan message) {
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