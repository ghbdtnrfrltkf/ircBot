package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
)

type Bot struct {
	nick    string
	user    string
	channel string
}

type Server struct {
	addr string
	port string
}

var NewBot = Bot{
	nick:    "simple_bot",
	channel: "#channel",
	user:    "simple_bot",
}

var NewServer = Server{
	addr: "127.0.0.1",
	port: "6667",
}

func (serv *Server) Connect() (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", serv.addr+":"+serv.port)
	if err != nil {
		log.Fatal("unable to connect to IRC server ", err)
	}

	log.Printf("Connected to IRC server %s (%s)\n", serv.addr, conn.RemoteAddr())
	return conn, nil
}

func main() {
	serv := NewServer
	ircBot := NewBot
	conn, _ := serv.Connect()
	conn.Write([]byte("USER " + ircBot.user + " 8 * :" + ircBot.user + "\r\n"))
	conn.Write([]byte("NICK " + ircBot.nick + "\r\n"))
	conn.Write([]byte("JOIN " + ircBot.channel + "\r\n"))

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	defer conn.Close()

	for {
		line, err := tp.ReadLine()
		if err != nil {
			break // break loop on errors
		}
		if strings.HasPrefix(line, "PING :") {
			PongKey := strings.Trim(line, "PING :")
			conn.Write([]byte("PONG :" + PongKey))
		}
		if strings.Contains(line, "REGISTER") {
			conn.Write([]byte("JOIN " + ircBot.channel + "\r\n"))
		}
		if strings.Contains(line, ":") {

			st := strings.Trim(line, "JOIN :#channel")
			ss := strings.Split(st, "!")[0]
			fmt.Fprintf(conn, "PRIVMSG "+ircBot.channel+" :Добро Пожаловать "+ss+" \r\n")
		}

		fmt.Printf("%s\n", line)
	}

}
