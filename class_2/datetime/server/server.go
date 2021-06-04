package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ln, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go sendMessage()
	go broadcaster(ctx)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}
		go handleConn(ctx, conn)
	}
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
			return
		}
	}
}

func handleConn(ctx context.Context, c net.Conn) {
	defer func() {
		if err := c.Close(); err != nil {
			return
		}
	}()

	ch := make(chan string)
	go clientWriter(c, ch)
	entering <- ch
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			ch <- "Current time: " + time.Now().Format("15:04:05")
			time.Sleep(1 * time.Second)
		case <-ctx.Done():
			ch <- "Server stopped"
			if err := c.Close(); err != nil {
				log.Print(err)
			}
		}
	}
}

func sendMessage() chan string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("send message: ")
		msg, _, err := reader.ReadLine()
		if err != nil {
			reader.Reset(os.Stdin)
			continue
		}
		messages <- "Message from server: " + string(msg)
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		if _, err := fmt.Fprintln(conn, msg); err != nil {
			log.Printf("error while print msg: %v", err)
		}
	}
}
