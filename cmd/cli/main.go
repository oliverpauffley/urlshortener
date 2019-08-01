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

	for {
		fmt.Println("\nUrl Shortener \n" +
			"To add a url to shorten type add: http:// \n" +
			"To redirect, type a short url: http:// \n" +
			"or type 'exit' to quit")
		cli.ReadInput()
	}

}
