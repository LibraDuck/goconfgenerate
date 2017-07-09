package main

import (
	"os"
	"io"
	"log"
)

func checkFileIsExist(filename string) bool {
	exist := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func Mkdir(path string) error {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir + path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func iowrite(w string) error {
	Mkdir("/config")
	dstFile, err := os.Create("config/config.go")
	if err != nil {
		return err
	}
	defer dstFile.Close()

	dstFile.WriteString(w)

	return err
}

func AppendToFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	} else {
		_, err := io.WriteString(f, content)
		if err != nil {
			log.Println(err)
		}
	}
	defer f.Close()
	return err
}

func Assembly() {
	var body string
	body += tplpackage
	body += tplimport
	iowrite(body)
}