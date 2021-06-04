package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	var nickname string
	fmt.Print("Please enter your nickname: ")
	if _, err = fmt.Scan(&nickname); err != nil {
		log.Fatal(err)
	}
	if _, err = fmt.Fprintf(conn, nickname); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = conn.Close(); err != nil {
			log.Printf("error while closing connection:%v", err)
		}
	}()

	go func() {
		_, err = io.Copy(os.Stdout, conn)
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: exit", conn.LocalAddr())
}
