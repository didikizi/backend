package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

var (
	addclient    = make(chan Profile)
	deleteclient = make(chan Profile)
	messagech    = make(chan string)
	answer       = make(chan Answer)
	result       int
	message      string
)

type Profile struct {
	name  string
	ch    chan string
	score int
}

type Answer struct {
	profile *Profile
	answer  string
}

func main() {
	listener, err := net.Listen("tcp", ":8001")
	message = ""
	if err != nil {
		panic(err)
	}
	go goWorker()
	go numbers()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go newclient(conn)
	}

}

func goWorker() {
	mapa := make(map[Profile]bool)
	for {
		select {
		case client := <-addclient:
			mapa[client] = true
			client.ch <- message
		case client := <-deleteclient:
			delete(mapa, client)
			close(client.ch)
		case msg := <-messagech:
			message = msg
			for client := range mapa {
				client.ch <- msg
			}
		}
	}
}

func newclient(conn net.Conn) {
	ch := make(chan string)
	client := Profile{"", ch, 0}
	go clientWriter(conn, client.ch)
	client.ch <- "Hi, what is your name?"
	scan := bufio.NewScanner(conn)
	for {
		scan.Scan()
		name := scan.Text()
		if len(name) > 3 {
			client.name = name
			break
		}
		client.ch <- "The name must consist of more than 3 characters"
	}
	addclient <- client
	for scan.Scan() {
		number := scan.Text()
		if len(number) >= 1 {
			answerClient := Answer{&client, number}
			answer <- answerClient
		}
	}
	deleteclient <- client
}

func numbers() {
	for {
		rand.Seed(time.Now().UnixNano())
		num1 := rand.Intn(100)
		num2 := rand.Intn(100)
		result = num1 + num2
		messagech <- strconv.Itoa(num1) + " + " + strconv.Itoa(num2)
		log.Println("Ansewr: ", result)
		for {
			msg := <-answer
			number, err := strconv.Atoi(msg.answer)
			if err == nil {
				if result == number {
					msg.profile.score++
					messagech <- msg.profile.name + " win!" + "\t Score:" + strconv.Itoa(msg.profile.score)
					break
				} else {
					msg.profile.ch <- "The answer is not correct"
				}
			} else {
				log.Println(err)
			}
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
