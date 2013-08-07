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
)

type Message struct {
	nickname    string
	message     string
	messageType MessageType
}

type Client struct {
	conn     net.Conn
	nickname string
	ch       chan Message
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
		if line == "/quit" {
			// QUIT
			io.WriteString(c.conn, "Bye!\n")
			c.conn.Close()
			ch <- Message{
				nickname:    c.nickname,
				message:     "",
				messageType: messageTypeQuit,
			}
		} else if strings.HasPrefix(line, "/me ") {
			// EMOTE
			ch <- Message{
				nickname:    c.nickname,
				message:     line[4:],
				messageType: messageTypeEmote,
			}
		} else {
			// SAY
			ch <- Message{
				nickname:    c.nickname,
				message:     line,
				messageType: messageTypeSay,
			}
		}
	}
}

func (c Client) WriteLinesFrom(ch <-chan Message) {
	for msg := range ch {
		toPrint := ""
		switch {
		case msg.messageType == messageTypeSay:
			toPrint = addColor(colorYellow, colorBlack, fmt.Sprintf("%s says %s", msg.nickname, msg.message))
		case msg.messageType == messageTypeEmote:
			toPrint = addColor(colorGreen, colorBlack, fmt.Sprintf("%s %s", msg.nickname, msg.message))
		case msg.messageType == messageTypeQuit:
			toPrint = addColor(colorRed, colorBlack, fmt.Sprintf("%s has quit.", msg.nickname))
		case msg.messageType == messageTypeJoin:
			toPrint = addColor(colorRed, colorBlack, fmt.Sprintf("%s has joined.", msg.nickname))
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
