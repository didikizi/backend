package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Profile struct {
	ch   chan string
	name string
}

var (
	entering = make(chan Profile)
	leaving  = make(chan Profile)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		panic(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[Profile]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)

		case cli := <-entering:
			clients[cli] = true
		}
	}
}

func handleConn(conn net.Conn) {
	tmp := make(chan string)
	profile := Profile{tmp, ""}
	go clientWriter(conn, profile.ch)

	who := conn.RemoteAddr().String()
	profile.ch <- "You are ip " + who + "\n Enter your name"
	input := bufio.NewScanner(conn)
	for {
		input.Scan()
		if len(input.Text()) > 4 {
			profile.name = input.Text()
			break
		}
		profile.ch <- "The name must be more than 4 letters"
	}

	entering <- profile
	messages <- profile.name + " has arrived"

	for input.Scan() {
		msg := profile.name + ": " + input.Text()
		if len(input.Text()) > 3 {
			messages <- msg
		} else {
			profile.ch <- "Message must be more than 3 characters"
		}
	}
	leaving <- profile
	messages <- profile.name + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
