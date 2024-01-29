package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	redisServerNetwork = "tcp"
	redisServerAddress = "0.0.0.0:6379"
)
const (
	simpleString = "+"
)

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

}

func run() error {

	listener, err := net.Listen(redisServerNetwork, redisServerAddress)

	if err != nil {
		return fmt.Errorf("failed to bind to port adress: %s : %v\n", redisServerAddress, err)
	}

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Printf("failed to close listener: %v\n", err)
			return
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connections: %v\n", err)
		}

		go handleClient(conn)
	}

}

func handleClient(conn net.Conn) {

	defer func(con net.Conn) {
		err := con.Close()
		if err != nil {
			log.Fatal("failed to close connection: ", err)
		}
	}(conn)

	for {
		handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {

	reader := bufio.NewReader(conn)
	line, err := readLine(reader)

	if err != nil {
		log.Printf("Error reading line: %v", err)
		return
	}

	var commands []string

	if strings.HasPrefix(line, "*") {
		n, err := strconv.Atoi(line[1:])
		if err != nil {
			log.Printf("Error converting string to int: %v", err)
			return
		}

		for i := 0; i < n; i++ {
			argLine, err := readLine(reader)

			if err != nil {
				log.Printf("Error reading argument line: %v", err)
				return
			}

			if strings.HasPrefix(argLine, "$") {
				argLen, _ := strconv.Atoi(argLine[1:])
				fmt.Printf("Length of argument: %d\n", argLen)

				n, _ := reader.ReadString('\n')
				out := strings.TrimSuffix(n, "\r\n")
				commands = append(commands, out)
				fmt.Printf("Arguments: %s\n", n)
			}
		}

	}

	res := handleCommands(commands)

	_, err := conn.Write(res)

}

func handleCommands(commands []string) []byte {
	var result []byte

	if len(commands) > 0 {
		switch commands[0] {
		case "echo":
			if len(commands) > 1 {
				fmt.Println(commands[1])
				result = simpleStringResponse(commands[1])
			}
		default:
			fmt.Println("Unknown command")
		}
	}

	return result
}

func simpleStringResponse(s string) []byte {
	result := simpleString + s + "\r\n"
	return []byte(result)
}

func readLine(reader *bufio.Reader) (string, error) {
	data, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(data, "\r\n"), nil
}

// TODO run yay -Rsu redis
