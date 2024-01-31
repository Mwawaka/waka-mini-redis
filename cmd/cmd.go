package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
}
