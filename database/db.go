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
		LongUrlFromHash(hash string) (string, error)
	}
	UrlModel struct {
		ID      int    `db:"id"`
		LongUrl string `db:"long_url"`
		Hash    string `db:"short_url"`
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
	hash := hashing.NewHashId(int(id))

	// add hashId into database
	sqlStatement = "UPDATE main.urls set hash =$1 WHERE id = $2"
	_, err = db.Exec(sqlStatement, hash, id)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (db *DB) HashFromLongUrl(longUrl string) (string, error) {
	sqlStatement := "SELECT hash FROM main.urls WHERE long_url = $1"
	row := db.QueryRow(sqlStatement, longUrl)
	var hash string

	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (db *DB) LongUrlFromHash(hash string) (string, error) {
	sqlStatement := "SELECT long_url FROM main.urls WHERE hash = $1"
	row := db.QueryRow(sqlStatement, hash)
	var longUrl string

	err := row.Scan(&longUrl)
	if err != nil {
		return "", err
	}

	return longUrl, nil
}
