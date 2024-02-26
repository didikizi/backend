package main

import (
	"context"
	"io"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	listener    net.Listener
	Connections chan net.Conn
}

func (s Server) Start() {
	log.Println("Server started on ", s.listener.Addr())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		s.Connections <- conn
	}
}

func main() {
	var wg sync.WaitGroup
	server := NewServer("8001")
	defer func() {
		err := server.listener.Close()
		Errpanic(err)
	}()

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	go server.Start()

	for {
		select {
		case <-ctx.Done():
			log.Println("start gracefull")
			wg.Wait()
			log.Println("stop gracefull")
			return
		case conn := <-server.Connections:
			wg.Add(1)
			go handleConn(ctx, conn, &wg)
		}
	}
}

func NewServer(port string) Server {
	lister, err := net.Listen("tcp", ":"+port)
	Errpanic(err)
	chn := make(chan net.Conn)
	return Server{
		listener:    lister,
		Connections: chn,
	}
}

func Errpanic(err error) {
	if err != nil {
		panic(err)
	}
}

func handleConn(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	t := time.NewTicker(time.Second)
	defer conn.Close()
	defer wg.Done()
	for {
		select {
		case time := <-t.C:
			io.WriteString(conn, time.Format("15:04:05\n"))
		case <-ctx.Done():
			io.WriteString(conn, "Bey!")
			return
		}
	}
}
