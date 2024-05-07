package database

import (
	"backend/models"
	"database/sql"
	"net/mail"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
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
gender TEXT,
age INTEGER DEFAULT 0,
 firstname TEXT,
 lastname TEXT
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
	r := "SELECT * FROM users WHERE email = ?;"
	row := db.QueryRow(r, c.Email.Address)
	err = row.Scan()
	if err != sql.ErrNoRows {
		return err
	}

	r = "INSERT INTO users VALUES(?, ?, ?, ?, ?, ?, ?, ?);"
	_, err = db.Exec(r,
		c.Uuid,
		c.Email.Address,
		c.Username,
		c.Password,
		c.Gender,
		c.Age,
		c.FirstName,
		c.LastName,
	)
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

func CreatePost(p models.Post) (err error) {
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		return err
	}
	defer db.Close()

	r := "INSERT INTO posts VALUES(?, ?, ?);"
	_, err = db.Exec(r, p.UserID, p.Content, time.Now())
	return err
}

func GetAllPosts() ([]models.Post, error) {
	res := []models.Post{}
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	r := "SELECT userid, users.username, content, date FROM posts JOIN users ON users.uuid = posts.userid;"
	rows, err := db.Query(r)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := models.Post{}
		err = rows.Scan(&p.UserID, &p.Username, &p.Content, &p.Date)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
}
