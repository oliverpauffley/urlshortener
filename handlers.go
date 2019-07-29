package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func (env *Env) urlHandler() http.HandlerFunc {
	type requestJson struct {
		LongUrl string `json:"long_url"`
	}
	type responseJson struct {
		ShortUrl string `json:"short_url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// check for POST or GET request

		if r.Method == "POST" {
			// decode json request
			newRequest := &requestJson{}
			err := json.NewDecoder(r.Body).Decode(newRequest)
			if err != nil {
				log.Printf("error decoding json request, %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// check url is valid
			resp, err := http.Get(newRequest.LongUrl)
			if err != nil {
				log.Printf("error when checking user url input, %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			} else if resp.StatusCode != http.StatusOK {
				log.Printf("problem with the url entered, check that the url to shorten is valid")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// check if url is already in database
			hashId, err := env.db.HashFromLongUrl(newRequest.LongUrl)
			if err == sql.ErrNoRows {
				// url is not in db so add
				hashId, err = env.db.NewUrl(newRequest.LongUrl)
			}
			if err != nil {
				log.Printf("error writing to database, %v", err)
			}

			// get url of server and add hashId to create a short url to return
			shortUrl := r.Host + "/" + hashId

			// encode and send json response
			var payload responseJson
			payload.ShortUrl = shortUrl
			b, err := json.Marshal(payload)
			if err != nil {
				log.Printf("error encoding json for response, %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, err = w.Write(b)
			if err != nil {
				log.Printf("error writing json to response body, %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// return short url in response
			return
		}
		// no other methods implemented so return error
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}
}
