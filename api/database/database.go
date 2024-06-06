package database

import (
	"database/sql"
	"io"
	"net/mail"
	"os"
	"real-time-forum/api/models"

	_ "github.com/mattn/go-sqlite3"
)

const (
	db_path = "api/data/database.sqlite"
	db_sql  = "api/data/db.sql"
)

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

func AddUser(u models.User) (err error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return err
	}
	defer db.Close()
	r := "SELECT * FROM users WHERE email = ?;"
	row := db.QueryRow(r, u.Email)
	err = row.Scan()
	if err != sql.ErrNoRows {
		return err
	}

	r = "INSERT INTO users VALUES(?, ?, ?, ?, ?, ?, ?, ?);"
	_, err = db.Exec(r,
		u.B64,
		u.Email,
		u.Name,
		u.GetPassword(),
		u.Gender,
		u.Age,
		u.FirstName,
		u.LastName,
	)
	return err
}

func GetPasswordAndIDFromMail(email mail.Address) (password []byte, b64 string, err error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, "", err
	}
	defer db.Close()

	r := "SELECT password, b64 FROM users WHERE email = ?;"
	row := db.QueryRow(r, email.Address)
	err = row.Scan(&password, &b64)
	if err != nil {
		return nil, "", err
	}

	return password, b64, err
}

func CreatePost(p models.Post) (err error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return err
	}
	defer db.Close()

	r := "INSERT INTO posts VALUES(NULL ,?, ?, ?);"

	_, err = db.Exec(r, p.UserID, p.Content, p.Date)
	return err
}

func GetAllPosts() ([]models.Post, error) {
	res := []models.Post{}
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	r := "SELECT user_id, users.name, content, date FROM posts JOIN users ON users.b64 = posts.user_id;"
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

	r := "SELECT b64, email, name, gender, age, first_name, last_name  FROM users;"
	rows, err := db.Query(r)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := models.User{}
		err = rows.Scan(&u.B64,
			&u.Email,
			&u.Name,
			&u.Gender,
			&u.Age,
			&u.FirstName,
			&u.LastName,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, u)
	}

	return res, nil
}
