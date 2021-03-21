package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
	"time"
)

type Bot struct {
	nick    string
	user    string
	channel string
	addr    string
	port    string
}

var NewBot = Bot{
	nick:    "WelcomeBot",
	channel: "#ru",
	user:    "WelcomeBot",
	addr:    "127.0.0.1",
	port:    "6667",
}

func (bot *Bot) Connect() (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", bot.addr+":"+bot.port)
	if err != nil {
		log.Fatal("unable to connect to IRC server ", err)
	}

	log.Printf("Connected to IRC server %s (%s)\n", bot.addr, conn.RemoteAddr())
	return conn, nil
}

func main() {
	ircBot := NewBot
	conn, _ := ircBot.Connect()
	conn.Write([]byte("USER " + ircBot.user + " 8 * :" + ircBot.user + "\r\n"))
	conn.Write([]byte("NICK " + ircBot.nick + "\r\n"))
	go func() {
		time.Sleep(15 * time.Second)
		conn.Write([]byte("JOIN " + ircBot.channel + "\r\n"))
	}()

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	defer conn.Close()

	for {
		line, err := tp.ReadLine()
		if err != nil {
			break // break loop on errors
		}
		if strings.HasPrefix(line, "PING :") {
			pongKey := strings.Trim(line, "PING :")
			conn.Write([]byte("PONG :" + pongKey + "\r\n"))
		}
		// if strings.Contains(line, "REGISTER") {
		// conn.Write([]byte("JOIN " + ircBot.channel + "\r\n"))
		// }
		if strings.Contains(line, "JOIN :#ru") {
			st := strings.Trim(line, ":")
			ss := strings.Split(st, "!")[0]
			if ss != "WelcomeBot" {
				fmt.Fprintf(conn, "PRIVMSG "+ircBot.channel+" :Добро Пожаловать "+ss+" \r\n")
			}
		}

		fmt.Printf("%s\n", line)
	}

}
