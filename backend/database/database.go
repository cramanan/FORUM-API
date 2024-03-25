package database

import (
	"backend/models"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const DB = "data/database.sqlite"

func InitDB() (err error) {
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		return err
	}
	defer db.Close()

	r := `CREATE TABLE IF NOT EXISTS clients (
id TEXT PRIMARY KEY,
email TEXT NOT NULL UNIQUE,
username TEXT,
password TEXT
);`
	_, err = db.Exec(r)
	return err
}

func AddClient(c models.Client) (err error) {
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		return err
	}
	defer db.Close()

	r := "INSERT INTO clients VALUES(?, ?, ?, ?)"

	_, err = db.Exec(r, c.Uuid, c.Email.Address, c.Username, c.Password)
	return err
}
