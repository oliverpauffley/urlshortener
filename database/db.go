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
	// open a db connection, will create a db with dataSource name if it doesn't exist.
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	// add table to db if it doesn't exist
	sqlStatement := "CREATE TABLE IF NOT EXISTS urls" +
		"(id INTEGER not null primary key unique," +
		"long_url TEXT    not null unique," +
		"hash     TEXT)"
	_, err = db.Exec(sqlStatement)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// NewUrl adds a url to db and generates a hash for the short url, uses transactions makes it safe for go routines
func (db *DB) NewUrl(longUrl string) (string, error) {
	sqlStatement, err := db.Prepare("INSERT INTO main.urls (long_url) VALUES ($1)")
	if err != nil {
		return "", err
	}
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	row, err := tx.Stmt(sqlStatement).Exec(longUrl)
	if err != nil {
		// error with adding to the db so rollback
		_ = tx.Rollback()
		return "", err
	} else {
		_ = tx.Commit()
	}

	// get id generated from SQLite
	id, err := row.LastInsertId()
	if err != nil {
		return "", nil
	}

	// generate hashid from integer id
	hash := hashing.NewHashId(int(id))

	// add hashId into database
	sqlStatement, err = db.Prepare("UPDATE main.urls set hash =$1 WHERE id = $2")
	if err != nil {
		return "", nil
	}

	tx, err = db.Begin()
	if err != nil {
		return "", nil
	}

	_, err = tx.Stmt(sqlStatement).Exec(hash, id)
	if err != nil {
		// error writing to db so rollback
		_ = tx.Rollback()
		return "", err
	} else {
		_ = tx.Commit()
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

// TearDown will delete tables from db after integration tests
func (db *DB) TearDown() error {
	_, err := db.Exec("DROP TABLE urls")
	if err != nil {
		return err
	}
	return nil
}
