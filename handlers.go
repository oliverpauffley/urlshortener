package urlshortener

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
			hashId, err := env.Db.HashFromLongUrl(newRequest.LongUrl)
			if err == sql.ErrNoRows {
				// url is not in db so add
				hashId, err = env.Db.NewUrl(newRequest.LongUrl)
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

func (env *Env) RedirectHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// check for POST or GET request
		if r.Method == "GET" {

			// get short url parameter from request url
			hashId := r.URL.Path[len("/"):]

			// get long url from db
			longUrl, err := env.Db.LongUrlFromHash(hashId)
			if err == sql.ErrNoRows {
				log.Printf("No url exists, check the url entered, %v", err)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if err != nil {
				log.Printf("error getting hash from database, %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// redirect to url
			http.Redirect(w, r, longUrl, http.StatusSeeOther)
			return
		}

		// no other request methods so return an error
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
