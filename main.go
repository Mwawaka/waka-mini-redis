package main

import (
	"log"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", "0.0.0.0:6370")

	if err != nil {
		log.Fatal("Failed to bind to port 6379", err)
	}

	defer func(lis net.Listener) {
		err := lis.Close()
		if err != nil {
			log.Fatal("Failed to close: ", err)
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting connections: ", err)
		}
		handleClient(conn)
	}

}

func handleClient(conn net.Conn) {
	defer func(con net.Conn) {
		err := con.Close()
		if err != nil {
			log.Fatal("Failed to close connection: ", err)
		}
	}(conn)
}
