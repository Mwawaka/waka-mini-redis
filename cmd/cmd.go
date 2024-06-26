package cmd

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	REDIS_RDB_OPCODE_AUX           = byte(0xFA)
	REDIS_RDB_OPCODE_RESIZEDB      = byte(0xFB)
	REDIS_RDB_OPCODE_EXPIRETIME_MS = byte(0xFC)
	REDIS_RDB_OPCODE_EXPIRETIME    = byte(0xFD)
	REDIS_RDB_OPCODE_SELECTDB      = byte(0xFE)
	REDIS_RDB_OPCODE_EOF           = byte(0xFF)
)

func Cmd() []string {
	//non var which returns a pointer that can be stored
	var args []string
	dir := flag.String("dir", " ", "directory where the RDB files are stored")
	filename := flag.String("dbfilename", "", "the name of the RDB file")
	flag.Parse()
	fmt.Println("directory: ", *dir)
	fmt.Println("database filename: ", *filename)

	args = append(args, *dir, *filename)
	createDirAndFile(*dir, *filename)
	return args
}

func createDirAndFile(dirname, filename string) {

	// get absolute path of the directory
	absDirPath, err := filepath.Abs(dirname)

	if err != nil {
		log.Fatal("error getting the absolute path of the directory")
	}

	// create directory if it doesn't exist
	if _, err := os.Stat(absDirPath); os.IsNotExist(err) {

		if err := os.MkdirAll(absDirPath, os.ModePerm); err != nil {
			log.Fatal("error creating the directory: ", err)
		}
	}

	// creating the file in the directory
	filePath := filepath.Join(absDirPath, filename)
	file, err := os.Create(filePath)

	if err != nil {
		log.Fatal("error creating the file")
	}

	if err = file.Close(); err != nil {
		log.Fatal("error closing the file")
	}
	fmt.Println(filePath)
	handleOpenFile()
}

func handleOpenFile() {
	const abs string = "/home/mwawaka/go/src/projects/mini-redis/dump.rdb"
	file, err := os.Open(abs)

	if err != nil {
		log.Fatal("error opening the file")
	}
	defer func(f *os.File) {
		err := f.Close()

		if err != nil {
			log.Fatal("error closing the file")
		}
	}(file)

	err = parseRDB(file)

	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal("error parsing: ", err)
	}

}

func parseRDB(file *os.File) error {
	reader := bufio.NewReader(file)
	magicString := make([]byte, 5)
	versionString := make([]byte, 4)
	_, err := io.ReadFull(reader, magicString)

	if err != nil {
		return err
	}

	_, err = io.ReadFull(reader, versionString)

	if err != nil {
		return err
	}

	for {
		var opcode byte
		err := binary.Read(reader, binary.LittleEndian, &opcode)

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		//fmt.Printf("%x, ", opcode)

		switch opcode {
		case REDIS_RDB_OPCODE_AUX:
			err = handleAux(reader)
			return err
		}

	}

	return nil
}

func handleAux(reader *bufio.Reader) error {
	keyLength, err := binary.ReadUvarint(reader)

	if err != nil {
		return err
	}
	fmt.Println(keyLength)
	keyBytes := make([]byte, keyLength)

	_, err = reader.Read(keyBytes)

	if err != nil {
		log.Print(err)
		return err
	}

	fmt.Printf("Key:%s\n", keyBytes)

	valueLength, err := binary.ReadUvarint(reader)

	if err != nil {
		return err
	}
	fmt.Println(valueLength)
	valueBytes := make([]byte, valueLength)

	_, err = reader.Read(valueBytes)

	if err != nil {
		log.Print(err)
		return err
	}

	fmt.Printf("Value:%s\n", valueBytes)
	return nil
}

//func handleString(reader *bufio.Reader) {
//	bit, _ := reader.ReadByte()
//	sb := bit >> 6
//	print(sb)
//	fmt.Printf("Byte:%x\n", bit)
//}
