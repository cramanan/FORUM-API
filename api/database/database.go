package database

import (
	"database/sql"
	"io"
	"os"
	"real-time-forum/api/models"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const db_path = "api/data/database.sqlite"
const db_sql = "api/data/db.sql"

func InitDB() (err error) {

	f, err := os.Open(db_sql)
	if err != nil {
		return err
	}
	defer f.Close()

	query, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(string(query))
	return err
}

func AddClient(c models.User) (err error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return err
	}
	defer db.Close()
	r := "SELECT * FROM users WHERE email = ?;"
	row := db.QueryRow(r, c.Email)
	err = row.Scan()
	if err != sql.ErrNoRows {
		return err
	}

	r = "INSERT INTO users VALUES(?, ?, ?, ?, ?, ?, ?, ?);"
	_, err = db.Exec(r,
		c.B64,
		c.Email,
		c.Username,
		c.Password,
		c.Gender,
		c.Age,
		c.FirstName,
		c.LastName,
	)
	return err
}

func GetClientFromMail(email string) (c *models.User, err error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	r := "SELECT b64, email, username, password FROM users WHERE email = ?;"
	row := db.QueryRow(r, email)
	c = new(models.User)
	err = row.Scan(&c.B64, &c.Email, &c.Username, &c.Password)
	if err != nil {
		return nil, err
	}

	return c, err
}

func CreatePost(p models.Post) (err error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return err
	}
	defer db.Close()

	r := "INSERT INTO posts VALUES(?, ?, ?, ?);"

	rawID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	_, err = db.Exec(r, rawID.String(), p.UserID, p.Content, time.Now())
	return err
}

func GetAllPosts() ([]models.Post, error) {
	res := []models.Post{}
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	r := "SELECT userid, users.username, content, date FROM posts JOIN users ON users.b64 = posts.userid;"
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

func GetAllUsers() ([]models.User, error) {
	res := []models.User{}
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	r := "SELECT username, b64 FROM users;"
	rows, err := db.Query(r)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := models.User{}
		err = rows.Scan(&u.Username, &u.B64)
		if err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}
