package main

import (
	"github.com/oliverpauffley/urlshortener"
	"github.com/oliverpauffley/urlshortener/database"
	"log"
	"net/http"
)

func main() {

	Db, err := database.NewDB("../../urlDB")
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()

	Router := http.NewServeMux()
	env := &urlshortener.Env{Db: Db, Router: Router}
	env.Routes()

	log.Fatal(http.ListenAndServe("localhost:8000", env.Router))
}
