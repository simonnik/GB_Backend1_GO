package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// read nickname from stdin
	fmt.Print("Enter your nickname:")
	reader := bufio.NewReader(os.Stdin)
	nickname, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("cannot read data, program exit")
		return
	}

	// establish connection
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// send nickname to server
	_, err = conn.Write(nickname)
	if err != nil {
		fmt.Printf("could not send nickname to server, %v", err)
		return
	}

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Print(err)
		}
	}()

	_, err = io.Copy(conn, os.Stdin) // until you send ^Z
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("%s: exit", conn.LocalAddr())
}
