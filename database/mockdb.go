package database

import (
	"database/sql"
	"github.com/oliverpauffley/urlshortener/hashing"
)

type Mockdb struct {
	Urls map[int]UrlModel
}

func (db Mockdb) NewUrl(longUrl string) (string, error) {

	// get maximum id from mockdb
	key := 0
	for range db.Urls {
		if _, exists := db.Urls[key]; exists == true {
			key++
		}
	}
	url := UrlModel{LongUrl: longUrl, ID: key, ShortUrl: hashing.NewHashId(key)}
	db.Urls[key] = url
	return url.ShortUrl, nil
}

func (db Mockdb) HashFromLongUrl(longUrl string) (string, error) {
	// iterate through the db for url
	for _, entry := range db.Urls {
		if entry.LongUrl == longUrl {
			// url is already in the db so return the short version
			return entry.ShortUrl, nil
		}
	}
	// url is not in the db so return an error
	return "", sql.ErrNoRows
}
