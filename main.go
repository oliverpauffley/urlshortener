package main

import (
	"github.com/oliverpauffley/urlshortener/database"
	"log"
	"net/http"
)

type Env struct {
	db     database.Store
	router *http.ServeMux
}

func main() {

	db, err := database.NewDB("./urlDB")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := http.NewServeMux()
	env := &Env{db, router}
	env.routes()

	log.Fatal(http.ListenAndServe("localhost:8000", env.router))
}
