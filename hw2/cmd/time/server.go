package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type client chan<- string

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
	go chs.sendMsgToClients()

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
				cli <- msg
			}

		case cli := <-s.entering:
			clients[cli] = true

		case cli := <-s.leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
func (s server) sendMsgToClients() {
	// канал сообщений для  клиентов
	userChan := make(chan string)
	// читаем сообщения из ввода сервера
	go func() {
		r := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Введите текст: ")
			text, _ := r.ReadString('\n')
			text = text[:len(text)-1] // обрезаем символ '\n'
			userChan <- text
		}
	}()

	ticker := time.NewTicker(time.Second)
	var message string
	for {
		// освобождаем канал тикера
		<-ticker.C
		select {
		case userMsg := <-userChan: // пришли пользовательские сообщения
			message = time.Now().Format("15:04:05") + " " + userMsg
		default: // никаких сообщений не пришло
			message = time.Now().Format("15:04:05")
		}
		s.messages <- message
	}
}

func (s server) handleConn(c net.Conn) {
	defer c.Close()
	ch := make(chan string)

	s.entering <- ch
	for msg := range ch {
		_, err := fmt.Fprintln(c, msg)
		if err != nil {
			break
		}
	}

	s.leaving <- ch
}
