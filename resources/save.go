package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	redisServerNetwork = "tcp"
	redisServerAddress = "0.0.0.0:6379"
)
const (
	simpleStrings  = "+"
	simpleErrors   = "-"
	integers       = ":"
	bulkStrings    = "$"
	arrays         = "*"
	nullBulkString = "-1"
	crlf           = "\r\n"
)

var (
	db = make(map[string]string)
)

func main() {

	err := cmd()
	err = run()

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
		if err := l.Close(); err != nil {
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

	if strings.HasPrefix(line, arrays) {
		n, _ := strconv.Atoi(line[1:])
		fmt.Printf("Length of array: %d\n", n)

		for i := 0; i < n; i++ {
			argLine, err := readLine(reader)

			if err != nil {
				log.Fatal("error reading string:", err)
			}

			if strings.HasPrefix(argLine, bulkStrings) {
				argLen, _ := strconv.Atoi(argLine[1:])
				fmt.Printf("Length of argument: %d\n", argLen)
				args, err := readLine(reader)
				fmt.Println(args)
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
			result = handleEcho(commands)
		case "PING":
			result = handlePing()
		case "GET":
			result = handleGet(commands)
		case "SET":
			result = handleSet(commands)
		default:
			result = simpleErrorResponse(commands[0])
		}
	}
	return result
}

func handlePing() []byte {
	return simpleStringResponse("PONG")
}

func handleEcho(commands []string) []byte {
	if len(commands) > 1 {
		echoString := strings.Join(commands[1:], " ")

		return bulkStringResponse(echoString)
	}
	return nullBulkStringResponse()
}

func handleGet(commands []string) []byte {
	key := commands[1]
	value := db[key]
	if value == "" {
		return nullBulkStringResponse()

	}
	return bulkStringResponse(value)
}

func handleSet(commands []string) []byte {
	var value string
	key := commands[1]
	value = strings.Join(commands[2:], "")

	if len(commands) > 4 {
		value = strings.Join(commands[2:len(commands)-2], " ")
		command := strings.ToUpper(commands[len(commands)-2])
		if command == "PX" {
			expiryMS, err := strconv.Atoi(commands[len(commands)-1])

			if err != nil {
				log.Fatal("Error parsing string")
			}
			db[key] = value
			timer := time.After(time.Duration(expiryMS) * time.Millisecond)
			go deleteKey(key, timer)
		} else {
			return simpleErrorResponse(command)
		}
	}
	db[key] = value

	return bulkStringResponse("OK")
}

func handleConfig() {

}
func simpleErrorResponse(msg string) []byte {
	err := fmt.Sprintf(": %s:  command not found", msg)
	result := simpleErrors + err + crlf
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

func nullBulkStringResponse() []byte {
	result := bulkStrings + nullBulkString + crlf
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

func deleteKey(key string, timer <-chan time.Time) {
	<-timer
	db[key] = ""
}

func cmd() error {
	//non var which returns a pointer that can be stored
	dir := flag.String("dir", " ", "directory where the RDB files are stored")
	filename := flag.String("filename", "", "the name of the RDB file")
	flag.Parse()

	if err := os.Mkdir(*dir, os.ModePerm); err != nil {
		return errors.New("error creating directory")
	}

	path, err := os.Getwd()

	if err != nil {
		return errors.New("error getting directory")
	}
	filePath := path + "/" + *dir + "/" + *filename
	file, err := os.Create(filePath)

	if err != nil {
		return errors.New("error creating file")
	}

	if err := file.Close(); err != nil {
		return errors.New("error closing file")
	}

	fmt.Println("directory: ", *dir)
	fmt.Println("database filename: ", *filename)
	return nil
}
