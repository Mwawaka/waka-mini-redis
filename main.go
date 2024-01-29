package main

import (
	"bufio"
	"fmt"
	"io"
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
	simpleStrings = "+"
	simpleErrors  = "-"
	integers      = ":"
	bulkStrings   = "$"
	arrays        = "*"
	crlf          = "\r\n"
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
		return fmt.Errorf("failed to bind to port adress: %s", redisServerAddress)
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
	var commands []string
	reader := bufio.NewReader(conn)
	line, err := readLine(reader)

	if err != nil {
		log.Fatal("error reading string:", err)
	}

	if strings.HasPrefix(line, "*") {
		n, _ := strconv.Atoi(line[1:])
		fmt.Printf("Length of array: %d\n", n)

		for i := 0; i < n; i++ {
			argLine, err := readLine(reader)

			if err != nil {
				log.Fatal("error reading string:", err)
			}

			if strings.HasPrefix(argLine, "$") {
				argLen, _ := strconv.Atoi(argLine[1:])
				fmt.Printf("Length of argument: %d\n", argLen)
				args, err := readLine(reader)

				if err != nil {
					log.Fatal("error reading string:", err)
				}

				commands = append(commands, args)

			}
		}

	}

	res := handleCommands(commands)
	_, err = conn.Write(res)

	if err != nil {
		log.Fatal("failed closing the connection: ", err)
	}

}

func handleCommands(commands []string) []byte {
	var result []byte

	if len(commands) > 0 {
		switch strings.ToUpper(commands[0]) {
		case "ECHO":
			if len(commands) > 1 {
				echoString := strings.Join(commands[1:], " ")

				result = bulkStringResponse(echoString)
			}
		case "PING":
			result = simpleStringResponse("PONG")
		default:
			err := fmt.Sprintf(": %s:  command not found", commands[0])
			result = simpleErrorResponse(err)
		}
	}
	return result
}

func simpleErrorResponse(msg string) []byte {
	result := simpleErrors + msg + crlf
	return []byte(result)
}

func simpleStringResponse(msg string) []byte {
	result := simpleStrings + msg + crlf
	return []byte(result)
}
func bulkStringResponse(msg string) []byte {
	size := strconv.Itoa(len(msg))
	result := bulkStrings + size + crlf + msg + crlf
	return []byte(result)
}

func readLine(reader *bufio.Reader) (string, error) {
	data, err := reader.ReadString('\n')

	if err != nil {
		if err == io.EOF {
			return "", nil
		}
		return "", err
	}

	return strings.TrimSuffix(data, "\r\n"), nil
}
