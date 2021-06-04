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

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error while closing conn: %v", err)
		}
	}()

	buf := make([]byte, 256)
	for {
		if _, err = conn.Read(buf); err == io.EOF {
			break
		}

		if _, err = io.WriteString(os.Stdout, fmt.Sprintf("Current time: %s", string(buf))); err != nil {
			log.Printf("error while writing string: %v", err)
		}
	}
}
