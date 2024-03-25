package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(name string) (err error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return err
	}
	defer db.Close()
	return err
}
