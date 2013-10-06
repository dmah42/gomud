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
  messageTypeTell
	messageTypeEmote
  messageTypeShout
  messageTypeJoin
  messageTypeQuit
	messageTypeWho
)

type message struct {
	from    client
  to      string
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
				from:    c,
				message:     "",
				messageType: messageTypeQuit,
			}
		// SAY
    case strings.HasPrefix(line, "/say "):
			ch <- message{
				from:    c,
				message:     line[5:],
				messageType: messageTypeSay,
			}
    // TELL
    case strings.HasPrefix(line, "/tell "):
      // TODO: use fields to parse better
      if player, err := players.get(line[6:]); err == nil {
        if conn, _ := player.isConnected(); conn {
          ch <- message{
            from: c,
            to: player.Nickname,
            message: line[6:],
            messageType: messageTypeTell,
          }
        } else {
          io.WriteString(c.conn, fmt.Sprintf("%q does not appear to be online.\n", line[6:]))
        }
      } else {
        io.WriteString(c.conn, fmt.Sprintf("%q.\n", err))
      }
		// EMOTE
		case strings.HasPrefix(line, "/me "):
			ch <- message{
				from:    c,
				message:     line[4:],
				messageType: messageTypeEmote,
			}
    // SHOUT
    case strings.HasPrefix(line, "/shout "):
      ch <- message{
        from: c,
        message: line[7:],
        messageType: messageTypeShout,
      }
		// WHO
		case line == "/who":
			io.WriteString(c.conn, setFg(colorWhite, fmt.Sprintf("%v\n", getConnected())))
			// FINGER
		case strings.HasPrefix(line, "/finger "):
			if player, err := players.get(line[8:]); err == nil {
				toPrint := setFg(colorWhite, fmt.Sprintf("%+v ", player.finger()))
				if c, _ := player.isConnected(); c {
					toPrint += setFgBold(colorGreen, "[online]\n")
				} else {
					toPrint += setFgBold(colorRed, "[offline]\n")
				}
				io.WriteString(c.conn, toPrint)
			} else {
				io.WriteString(c.conn, fmt.Sprintf("%q.\n", err))
			}
		case line == "look":
			room, err := rooms.get(c.player.Room)
			if err == nil {
				io.WriteString(c.conn, room.describe())
			} else {
				// TODO: handle limbo
				io.WriteString(c.conn, fmt.Sprintf("%q.\n", err))
				log.Printf("%q in limbo %q.\n", c.player.Nickname, c.player.Room)
			}
		default:
      // TODO: handle exits, etc
      log.Printf("Unknown command: %q\n", line)
		}
	}
}

func sameRoom(c client, msg message) bool {
  return c.player.Room == msg.from.player.Room
}

func (c client) writeLinesFrom(ch <-chan message) {
	for msg := range ch {
    from := msg.from.player.Nickname
		toPrint := ""
		// TODO: Register command per message type for colors/format string and location restriction
		switch msg.messageType {
		case messageTypeSay:
      if sameRoom(c, msg) {
			  toPrint = setFg(colorYellow, fmt.Sprintf("%s says %s", from, msg.message))
      }
    case messageTypeTell:
      if msg.to == c.player.Nickname {
        toPrint = setFg(colorGreen, fmt.Sprintf("%s tells you %s", from, msg.message))
      }
		case messageTypeEmote:
      if sameRoom(c, msg) {
        toPrint = setFg(colorMagenta, fmt.Sprintf("%s %s", from, msg.message))
      }
    case messageTypeShout:
      toPrint = setFgBold(colorCyan, fmt.Sprintf("%s shouts %s", from, msg.message))
		case messageTypeQuit:
      if sameRoom(c, msg) {
			  toPrint = setFgBold(colorRed, fmt.Sprintf("%s has quit.", from))
      }
		case messageTypeJoin:
      if sameRoom(c, msg) {
        toPrint = setFgBold(colorRed, fmt.Sprintf("%s has joined.", from))
      }
		default:
			log.Printf("Unhandled message type: %+v", msg)
			continue
		}
    if len(toPrint) == 0 {
      continue
    }
		_, err := io.WriteString(c.conn, toPrint+"\n")
		if err != nil {
      log.Printf("Error writing '%q'\n", toPrint)
		}
	}
}
