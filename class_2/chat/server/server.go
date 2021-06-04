package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	ln, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go broadcaster(ctx)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error while connecting: %v", err)
		}
		go handleConn(ctx, conn)
	}
	cancel()
}

func broadcaster(ctx context.Context) {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		case <-ctx.Done():
			for cli := range clients {
				cli <- "server stopped"
			}
			return
		}
	}
}

func handleConn(ctx context.Context, conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}

	select {
	case leaving <- ch:
		messages <- who + " has left"
		if err := conn.Close(); err != nil {
			log.Printf("error while closing conn: %v", err)
			return
		}
	case <-ctx.Done():
		messages <- "server stopped"
		return
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			log.Printf("error while writting message for clients: %v", err)
		}
	}
}
