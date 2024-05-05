package database

import (
	"backend/models"
	"database/sql"
	"net/mail"

	"github.com/gofrs/uuid"
	"github.com/mattn/go-sqlite3"
)

const DB = "data/database.sqlite"

func InitDB() (err error) {
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		return err
	}
	defer db.Close()

	r := `CREATE TABLE IF NOT EXISTS users (
uuid TEXT PRIMARY KEY,
email TEXT NOT NULL UNIQUE,
username TEXT,
password TEXT,
gender TEXT
);

CREATE TABLE IF NOT EXISTS posts (
	userid TEXT REFERENCES users(uuid),
	content TEXT
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
	r := "SELECT * FROM users WHERE email = ?"
	row := db.QueryRow(r, c.Email.Address)
	if row.Err() == nil {
		return sqlite3.ErrConstraintUnique
	}

	r = "INSERT INTO users VALUES(?, ?, ?, ?, ?)"
	_, err = db.Exec(r, c.Uuid, c.Email.Address, c.Username, c.Password, c.Gender)
	return err
}

func GetClientFromMail(email *mail.Address) (c *models.Client, err error) {
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	r := "SELECT uuid, email, username, password FROM users WHERE email = ?;"
	row := db.QueryRow(r, email.Address)
	var id, str_mail string
	c = &models.Client{}
	err = row.Scan(&id, &str_mail, &c.Username, &c.Password)
	if err != nil {
		return nil, err
	}

	c.Uuid, err = uuid.FromString(id)
	if err != nil {
		return nil, err
	}

	c.Email, err = mail.ParseAddress(str_mail)
	if err != nil {
		return nil, err
	}
	return c, err
}

func GetAllPosts() (res []models.Post, err error) {
	res = []models.Post{}
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	r := "SELECT userid, users.username, content FROM posts JOIN users ON users.uuid = posts.userid;"
	rows, err := db.Query(r)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := models.Post{}
		err = rows.Scan(&p.UserID, &p.Username, &p.Content)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
}
