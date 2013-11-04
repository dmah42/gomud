package mudlib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

const logFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile

var msgchan = make(chan message)

var config = &struct {
	Port        int
	LoginPrompt string
}{}

var errorLogFile, statsLogFile os.File
var errorLog, statsLog *log.Logger

func init() {
	os.MkdirAll("logs/", os.ModePerm)
	errorLogFile, err := os.OpenFile("logs/error", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create error log: %+v", err)
		return
	}
	// Seek to end of file
	errorLogFile.Seek(0, 2)
	errorLog = log.New(errorLogFile, "[error] ", logFlags)

	statsLogFile, err := os.OpenFile("logs/stats", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create stats log: %+v", err)
		return
	}
	// Seek to end of file
	statsLogFile.Seek(0, 2)
	statsLog = log.New(statsLogFile, "[stats] ", logFlags)

	c := time.Tick(1 * time.Minute)
	go writeStats(c)
}

// Run listens for connections and handles player interaction.
func Run(configFile string) error {
	b, err := loadBytes(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, config)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(config.Port))
	if err != nil {
		return err
	}
	log.Printf("TCP listening on " + strconv.Itoa(config.Port))

	addchan := make(chan client)
	rmchan := make(chan client)

	go handleMessages(addchan, rmchan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn, addchan, rmchan)
	}

	if err := errorLogFile.Close(); err != nil {
		log.Printf("Failed to close error log: %+v", err)
		return err
	}
	if err := statsLogFile.Close(); err != nil {
		log.Printf("Failed to close stats log: %+v", err)
		return err
	}
	return nil
}

func writeStats(c <-chan time.Time) {
	memStats := new(runtime.MemStats)
	gcStats := new(debug.GCStats)
	debug.ReadGCStats(gcStats)
	statsLog.Printf("GC: %+v", *gcStats)
	runtime.ReadMemStats(memStats)
	statsLog.Printf("MEM: %+v", *memStats)
}

func handleConnection(c net.Conn, addchan chan<- client, rmchan chan<- client) {
	bufc := bufio.NewReader(c)
	defer c.Close()

	nick := promptNick(c, bufc)
	if nick == nil {
		io.WriteString(c, "Goodbye.")
		return
	}
	os.MkdirAll("logs/player/", os.ModePerm)
	logFile, err := os.OpenFile("logs/player/"+*nick, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		io.WriteString(c, "Something went wrong.\n")
		errorLog.Printf("Failed to create player log %q", logFile)
		return
	}

	newClient := client{
		conn:   c,
		player: *nick,
		ch:     make(chan message),
		log:    log.New(logFile, "["+*nick+"] ", logFlags),
	}

	addchan <- newClient
	msgchan <- message{
		from:        newClient,
		message:     "",
		messageType: messageTypeJoin,
	}
	defer func() {
		msgchan <- message{
			from:        newClient,
			message:     "",
			messageType: messageTypeQuit,
		}
		player, err := players.get(newClient.player)
		if err != nil {
			errorLog.Fatalf("%+v", err)
		}
		player.disconnect()
		newClient.log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmchan <- newClient
		if err := logFile.Close(); err != nil {
			errorLog.Printf("Failed to close player logFile %q", newClient.player)
		}
	}()

	player, err := players.get(newClient.player)
	if err != nil {
		errorLog.Fatalf("%+v", err)
	}
	player.connect()

	// Add player to room
	if err := player.toRoom(newClient, player.room); err != nil {
		errorLog.Printf("User %q is starting in unknown room %q\n", player.nickname, player.room)
		// TODO: handle limbo
		return
	}

	// Startup commands
	// TODO: allow player to set these
	if err := doCommand(newClient, "/look", []string{}); err != nil {
		errorLog.Printf("Failed to 'look' on startup\n")
	}

	go newClient.readLines()
	newClient.writeLinesFrom(newClient.ch)
}

func handleMessages(addchan <-chan client, rmchan <-chan client) {
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
			c.log.Printf("Client disconnect: %v\n", c.conn)
			log.Printf("Client disconnect: %v\n", c.conn)
			delete(clients, c.conn)
		}
	}
}

func promptNick(c net.Conn, bufc *bufio.Reader) (nick *string) {
	// TODO: custom prompts
	io.WriteString(c, setFgBold(colorMagenta, config.LoginPrompt)+"\n")
	errorCount := 0
	var realname string
	nick = new(string)
	for errorCount < 3 {
		io.WriteString(c, "What is your nick? ")
		nickBytes, _, _ := bufc.ReadLine()
		*nick = string(nickBytes)
		// TODO: password
		// Check for existing player.
		player, err := players.get(*nick)
		if err == nil {
			// check if user is already logged in
			if con, _ := player.isConnected(); con {
				io.WriteString(c, setFgBold(colorRed, fmt.Sprintf("%s is already connected. Please try again.\n", *nick)))
				errorCount++
				continue
			}

			io.WriteString(c, setFgBold(colorGreen, fmt.Sprintf("Welcome back, %s!\n", *nick)))
			log.Printf("Player %+v logged in.\n", player)
			return
		}
		// Not found so create a new one.
		log.Printf("Creating new player: %s\n", *nick)
		io.WriteString(c, "You seem to be new here. What is your real name? ")
		realnameBytes, _, _ := bufc.ReadLine()
		realname = string(realnameBytes)
		log.Printf("Adding real name %s for %s\n", realname, *nick)
		if player, err = players.add(*nick, realname); err == nil {
			io.WriteString(c, setFgBold(colorGreen, fmt.Sprintf("Welcome, %s!\n", *nick)))
			player.room = startRoomId
			return
		}
		errorLog.Printf("Error creating new player %s %s: %+v\n", *nick, realname, err)
		errorCount++
	}
	return
}
