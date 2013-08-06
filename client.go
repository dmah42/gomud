package gomud

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
		if strings.HasPrefix(line, "/me ") {
			ch <- Message{
				nickname:    c.nickname,
				message:     line[4:-1],
				messageType: messageTypeEmote,
			}
		} else {
			ch <- Message{
				nickname:    c.nickname,
				message:     line[:-1],
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
			toPrint = fmt.Sprintf("\033[1;33;40m%s says %s\033[m", msg.nickname, msg.message)
		case msg.messageType == messageTypeEmote:
			toPrint = fmt.Sprintf("\033[1;32;40m%s %s\033[m", msg.nickname, msg.message)
		case msg.messageType == messageTypeQuit:
			toPrint = fmt.Sprintf("\033[1;31;40m%s has quit.\033[m", msg.nickname)
		case msg.messageType == messageTypeJoin:
			toPrint = fmt.Sprintf("\033[1;30;40m%s has joined.\033[m", msg.nickname)
		default:
			log.Printf("Unknown message type: %+v", msg)
			return
		}
		_, err := io.WriteString(c.conn, toPrint)
		if err != nil {
			return
		}
	}
}
