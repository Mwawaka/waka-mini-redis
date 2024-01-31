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
	filename := flag.String("filename", "", "the name of the RDB file")
	flag.Parse()

	path, err := os.Getwd()

	if err != nil {
		log.Fatal("error getting directory")
	}

	dirPath := filepath.Join(path, *dir)

	_, err = os.Stat(dirPath)

	if os.IsNotExist(err) {
		if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
			log.Fatal("error creating directory: ", err)
		}
	}
	filePath := filepath.Join(dirPath, *filename)
	file, err := os.Create(filePath)

	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}

	if err := file.Close(); err != nil {
		log.Fatal("error closing file")
	}

	fmt.Println("directory: ", *dir)
	fmt.Println("database filename: ", *filename)
	args = append(args, *dir, *filename)
	return args
}
