package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	fileP string
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
	fileP = filePath
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
	var array []byte
	reader := bufio.NewReader(file)

	for {
		opCode, err := reader.ReadByte()

		if err != nil {
			return err
		}

		fmt.Print(opCode, ", ")
		array = append(array, opCode)

	}

}
