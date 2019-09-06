package db

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
)

type DB struct {
	db *dbx.DB
}

func New(link string) (*DB, error) {
	db, err := dbx.Open("postgres", link)
	return &DB{db: db}, err
}
