package main

import (
	"fmt"
	"log"
	"net"
)

const (
	addressString = "0.0.0.0:6379"
)

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

}

func run() error {
	listener, err := net.Listen("tcp", addressString)

	if err != nil {
		return fmt.Errorf("failed to bind to port adress: %s : %v\n", addressString, err)
	}

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Printf("failed to close listener: %v\n", err)
		}
	}(listener)
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connections: %v\n", err)
		}
		handleClient(conn)
	}

}

func handleClient(conn net.Conn) {
	defer func(con net.Conn) {
		err := con.Close()
		if err != nil {
			log.Fatal("failed to close connection: ", err)
		}
	}(conn)

	//Read Data
	for {
		buf := make([]byte, 128)

		n, err := conn.Read(buf)

		if err != nil {
			log.Println("error reading from buffer: ", err)
			return
		}
		log.Printf("data read: %s\n", buf[:n])

		res := "+PONG\r\n"
		_, err = conn.Write([]byte(res))

		if err != nil {
			log.Println("error writing to buffer: ", err)
		}
	}

}
