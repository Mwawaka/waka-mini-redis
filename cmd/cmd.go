package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func Cmd() error {
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
