package main

import (
	"bytes"
	"encoding/json"
	"github.com/oliverpauffley/urlshortener/database"
	"github.com/oliverpauffley/urlshortener/hashing"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUrlHandler(t *testing.T) {

	// set up mockdb and env for testing
	db := database.Mockdb{Urls: map[int]database.UrlModel{}}
	db.Urls[1] = database.UrlModel{LongUrl: "http://bbc.com", ID: 1, Hash: hashing.NewHashId(1)}
	env := Env{db, http.NewServeMux()}

	var tt = []struct {
		name     string // name of test
		inputUrl string // long url to shorten
		wantCode int    // http status code expected
		wantUrl  string // shortened url
	}{
		{
			name:     "returns bad request when sent an invalid url",
			inputUrl: "http://invalidurl.thisshouldntwork",
			wantCode: http.StatusBadRequest,
			wantUrl:  "",
		},
		{
			name:     "returns correct short url from valid long url",
			inputUrl: "http://google.com",
			wantCode: http.StatusOK,
			wantUrl:  "/" + hashing.NewHashId(0),
		},
		{
			name:     "returns short url when given a url already in the db",
			inputUrl: "http://bbc.com",
			wantCode: http.StatusOK,
			wantUrl:  "/" + hashing.NewHashId(1),
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			// format url as json
			requestPayload := struct {
				LongUrl string `json:"long_url"`
			}{
				LongUrl: test.inputUrl,
			}
			b, _ := json.Marshal(requestPayload)

			request, _ := http.NewRequest("POST", "/url", bytes.NewBuffer(b))
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(env.urlHandler())
			handler.ServeHTTP(response, request)

			gotCode := response.Code
			if gotCode != test.wantCode {
				t.Errorf("Got incorrect response code, wanted %v, got %v", test.wantCode, gotCode)
			}

			// decode json response if the test expects a url to be returned
			if test.wantUrl != "" {
				type responsePayload struct {
					ShortUrl string `json:"short_url"`
				}
				gotJson := &responsePayload{}
				err := json.NewDecoder(response.Body).Decode(gotJson)
				if err != nil {
					t.Errorf("error decoding json response, %v", err)
				}

				if gotJson.ShortUrl != test.wantUrl {
					t.Errorf("server returned the incorrect shortened url, wanted %v, got %v",
						test.wantUrl, gotJson.ShortUrl)
				}
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	// set up mockdb and env for testing
	db := database.Mockdb{Urls: map[int]database.UrlModel{}}
	db.Urls[1] = database.UrlModel{LongUrl: "http://bbc.com", ID: 1, Hash: hashing.NewHashId(1)}
	env := Env{db, http.NewServeMux()}

	var tt = []struct {
		name     string // name of test
		inputUrl string // short url to redirect
		wantCode int    // http status code expected
		wantUrl  string // redirect url
	}{
		{
			name:     "returns a redirect to a url when given a short url from the db",
			inputUrl: db.Urls[1].Hash,
			wantCode: http.StatusSeeOther,
			wantUrl:  "http://bbc.com",
		},
		{
			name:     "returns an error when trying to use a short url not in the db",
			inputUrl: "this isn't in the db",
			wantCode: http.StatusNotFound,
			wantUrl:  "",
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {

			request, _ := http.NewRequest("GET", "/"+test.inputUrl, nil)
			response := httptest.NewRecorder()

			handler := http.HandlerFunc(env.RedirectHandler())
			handler.ServeHTTP(response, request)

			gotCode := response.Code
			if gotCode != test.wantCode {
				t.Errorf("Got incorrect response code, wanted %v, got %v", test.wantCode, gotCode)
			}

			if test.wantUrl != "" {
				redirectUrl := response.Header().Get("Location")

				if redirectUrl != test.wantUrl {
					t.Errorf("server returned the incorrect shortened url, wanted %v, got %v",
						test.wantUrl, redirectUrl)
				}
			}
		})
	}
}
