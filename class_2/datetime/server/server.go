package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ln, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}
		go handleConn(conn, ctx)
	}
}

func handleConn(c net.Conn, ctx context.Context) {
	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("error while closing: %v", err)
		}
	}()

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ctx.Done():
			if _, err := io.WriteString(c, "Server was stopped"); err != nil {
				log.Printf("error while stopping server: %v", err)
			}
		case <-ticker.C:
			if _, err := io.WriteString(c, time.Now().Format("15:04:05\n\r")); err != nil {
				return
			}
		}
	}
}
