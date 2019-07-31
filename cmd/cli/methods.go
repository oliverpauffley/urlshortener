package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"github.com/oliverpauffley/urlshortener/database"
	"github.com/pkg/browser"
	"io"
	"net/http"
	"os"
	"strings"
)

type CLI struct {
	db database.Store
	in io.Reader
}

// ReadInput takes user input and adds urls and redirects if user requests it
func (cli *CLI) ReadInput() {
	reader := bufio.NewScanner(cli.in)
	reader.Scan()

	userInput := reader.Text()

	// split user input on the space and then use a switch to decide what to do
	parts := strings.Split(userInput, " ")
	if len(parts) > 2 {
		println("Usage: add: or url: followed by a space and your url")
		cli.ReadInput()
	}

	var err error
	switch parts[0] {
	case "add:":
		shortUrl, err := cli.AddUrl(parts[1])
		if err != nil {
			fmt.Printf("error adding url to database, %v", err)
			return
		}
		println("Short url created at:")
		println(shortUrl)

	case "url:":
		err = cli.Redirect(parts[1])
		if err != nil {
			fmt.Printf("error redirecting to url, %v", err)
			return
		}

	case "exit":
		os.Exit(1)

	default:
		println("Usage: add: or url: followed by a space and your url")
		cli.ReadInput()
	}
}

// AddUrl adds user inputted long url to the database
func (cli *CLI) AddUrl(longUrl string) (string, error) {
	// check url is valid
	resp, err := http.Get(longUrl)
	if err != nil {
		return "", err
	} else if resp.StatusCode != http.StatusOK {
		return "", errors.New("problem with the url you have entered. Check that the url you have entered is correct")
	}

	// check if url is already in the db
	hash, err := cli.db.HashFromLongUrl(longUrl)
	if err == sql.ErrNoRows {
		hash, err = cli.db.NewUrl(longUrl)
	}
	if err != nil {
		return "", err
	}

	// return short url, hardcoded url at the moment not sure how you would do this?
	shortUrl := "localhost:8000/" + hash

	return shortUrl, nil

}

func (cli *CLI) Redirect(shortUrl string) error {
	// separate hash from shorturl
	hash := shortUrl[len("http://localhost:8000/"):]

	// check if hash is in database
	longUrl, err := cli.db.LongUrlFromHash(hash)
	if err == sql.ErrNoRows {
		err = errors.New("no long url exists for this short url")
		return err
	}
	if err != nil {
		return err
	}

	// tell client to open browser with the url
	err = browser.OpenURL(longUrl)
	if err != nil {
		return err
	}
	return nil
}
