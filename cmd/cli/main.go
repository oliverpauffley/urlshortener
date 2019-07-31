package main

import (
	"fmt"
	"github.com/oliverpauffley/urlshortener/database"
	"log"
	"os"
)

func main() {
	db, err := database.NewDB("../../urlDB")
	if err != nil {
		log.Fatalf("error creating or opening database, %v", err)
	}
	cli := &CLI{db, os.Stdin}

	fmt.Println("Url Shortener \n" +
		"To add a url to shorten type add: http:// \n" +
		"To enter a short url type url: http:// \n" +
		"or type `exit` to quit")
	cli.ReadInput()

	os.Exit(0)
}
