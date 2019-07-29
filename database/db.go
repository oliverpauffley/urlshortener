package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oliverpauffley/urlshortener/hashing"
)

// interface for url db. mocked for unit tests in mockdb
type (
	Store interface {
		NewUrl(longUrl string) (string, error)
		HashFromLongUrl(longUrl string) (string, error)
	}
	UrlModel struct {
		ID       int    `db:"id"`
		LongUrl  string `db:"long_url"`
		ShortUrl string `db:"short_url"`
	}
)

type DB struct {
	*sql.DB
}

// initializes a new db
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) NewUrl(longUrl string) (string, error) {
	sqlStatement := "INSERT INTO main.urls (long_url) VALUES ($1)"
	row, err := db.Exec(sqlStatement, longUrl)
	if err != nil {
		return "", err
	}

	// get id generated from SQLite
	id, err := row.LastInsertId()
	if err != nil {
		return "", nil
	}

	// generate hashid from integer id
	hashId := hashing.NewHashId(int(id))

	// add hashId into database
	sqlStatement = "UPDATE main.urls set short_url =$1 WHERE id = $2"
	_, err = db.Exec(sqlStatement, hashId, id)
	if err != nil {
		return "", err
	}

	return hashId, nil
}

func (db *DB) HashFromLongUrl(longUrl string) (string, error) {
	sqlStatement := "SELECT short_url FROM main.urls WHERE long_url = $1"
	row := db.QueryRow(sqlStatement, longUrl)
	var hashId string

	err := row.Scan(&hashId)
	if err != nil {
		return "", err
	}

	return hashId, nil
}
