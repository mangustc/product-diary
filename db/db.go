package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	DB        *sql.DB
	TableName string
}

func NewStore(dbName string, tableName string, createQuery string) (*Store, error) {
	db, err := getDB(dbName)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database (%s)", err.Error())
	}

	if err := newTable(db, createQuery); err != nil {
		return nil, fmt.Errorf("Failed to create table (%s)", err.Error())
	}

	return &Store{
		DB:        db,
		TableName: tableName,
	}, nil
}

func getDB(dbName string) (*sql.DB, error) {
	// Init SQLite3 database
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func newTable(db *sql.DB, createQuery string) error {
	_, err := db.Exec(createQuery)
	if err != nil {
		return err
	}

	return nil
}
