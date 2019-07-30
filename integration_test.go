package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/oliverpauffley/urlshortener/database"
	"github.com/oliverpauffley/urlshortener/hashing"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var env *Env

// integration tests will only run when not in short mode using flags
func TestMain(m *testing.M) {
	flag.Parse()

	if !testing.Short() {
		// creating a sqlite db for testing
		db, err := database.NewDB("./testDB")
		if err != nil {
			log.Fatalf("error creating db for test, %v", err)
		}
		router := http.NewServeMux()
		env = &Env{db, router}
		env.routes()

		result := m.Run()

		// teardown and delete testing db
		err = db.TearDown()
		if err != nil {
			log.Fatalf("error tearing down test db, %v", err)
		}
		err = os.Remove("./testDB")
		if err != nil {
			log.Fatalf("error deleting test db. %v", err)
		}
		os.Exit(result)
	}
}

func TestIntegration(t *testing.T) {

	t.Run("a working url can be added to the db and then redirected to", func(t *testing.T) {
		testUrl := "http://google.com"
		// create json request
		requestPayload := struct {
			LongUrl string `json:"long_url"`
		}{
			LongUrl: testUrl,
		}
		b, err := json.Marshal(requestPayload)
		if err != nil {
			t.Fatalf("error encoding json, %v", err)
		}

		req, err := http.NewRequest("POST", "/url", bytes.NewBuffer(b))
		if err != nil {
			t.Fatalf("error creating new request, %v", err)
		}
		response := httptest.NewRecorder()

		env.router.ServeHTTP(response, req)

		// check http response code is correct
		gotCode := response.Code
		if gotCode != http.StatusOK {
			t.Errorf("got the incorrect response code when adding a url, wanted %v, got %v", http.StatusOK, gotCode)
		}

		// decoding json response from server and then check it is correct
		var responsePayload struct {
			ShortUrl string `json:"short_url"`
		}
		err = json.Unmarshal(response.Body.Bytes(), &responsePayload)
		if err != nil {
			t.Fatalf("error decoding json from urlhandler response, %v", err)
		}
		gotUrl := responsePayload.ShortUrl
		wantUrl := "/" + hashing.NewHashId(1)

		if gotUrl != wantUrl {
			t.Errorf("got incorrect short url from db, wanted %s, got %s", wantUrl, gotUrl)
		}

		// form new request using shortUrl from response
		req, err = http.NewRequest("GET", gotUrl, nil)
		if err != nil {
			t.Errorf("error creating new request, %v", err)
		}
		response = httptest.NewRecorder()

		env.router.ServeHTTP(response, req)

		// check for correct response code on redirect
		gotCode = response.Code
		if gotCode != http.StatusSeeOther {
			t.Errorf("got incorrect redirect code, wanted %v, got %v", http.StatusSeeOther, gotCode)
		}

		// check for correct url on redirect.
		gotUrl = response.Header().Get("Location")
		if gotUrl != testUrl {
			t.Errorf("got incorrect url on redirect, wanted %v, got %v", testUrl, gotUrl)
		}
	})
}
