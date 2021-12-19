package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client struct {
	ch   chan<- string
	nick string
}

type server struct {
	entering chan client
	leaving  chan client
	messages chan string
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	chs := server{
		entering: make(chan client),
		leaving:  make(chan client),
		messages: make(chan string),
	}

	go chs.broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go chs.handleConn(conn)
	}
}

func (s server) broadcaster() {
	clients := make(map[client]bool)

	for {
		select {
		case msg := <-s.messages:
			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-s.entering:
			clients[cli] = true

		case cli := <-s.leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

func (s server) handleConn(c net.Conn) {
	defer c.Close()
	ch := make(chan string)

	go s.clientWriter(c, ch)
	buf := make([]byte, 100) // создаем буфер
	_, err := c.Read(buf)
	if err != nil {
		fmt.Printf("could not read nickname from %s", c.RemoteAddr().String())
		return
	}
	who := string(buf)

	ch <- "your nickname is set to " + who
	s.messages <- who + " has arrived"
	s.entering <- client{ch, who}

	inputMsg := bufio.NewScanner(c)
	for inputMsg.Scan() {
		s.messages <- who + ": " + inputMsg.Text()
	}
	s.leaving <- client{ch, who}
	s.messages <- who + " has left"
}

func (s server) clientWriter(c net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := fmt.Fprintln(c, msg)
		if err != nil {
			break
		}
	}
}
