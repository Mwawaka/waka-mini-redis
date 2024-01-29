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

	//for {
	//
	//	buf := make([]byte, 128)
	//
	//	n, err := conn.Read(buf)
	//
	//	if err != nil {
	//		log.Println("error reading from buffer: ", err)
	//		return
	//	}
	//	log.Printf("data read: %s", string(buf[:n]))
	//
	//	res := "+PONG\r\n"
	//	_, err = conn.Write([]byte(res))
	//
	//	if err != nil {
	//		log.Println("error writing to buffer: ", err)
	//	}
	//}
}

func handleRequest(conn net.Conn) {

	reader := bufio.NewReader(conn)

	data, _ := reader.ReadString('\n')
	line := strings.TrimSuffix(data, "\r\n")
	var commands []string

	if strings.HasPrefix(line, "*") {
		n, _ := strconv.Atoi(line[1:])
		fmt.Printf("Length of array: %d\n", n)

		for i := 0; i < n; i++ {
			argLine, _ := reader.ReadString('\n')
			argLine = strings.TrimSuffix(argLine, "\r\n")

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
