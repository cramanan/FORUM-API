package database

import (
	"database/sql"
	"io"
	"os"
	"real-time-forum/api/models"

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
		u.Password,
		u.Gender,
		u.Age,
		u.FirstName,
		u.LastName,
	)
	return err
}

func GetUserFromMail(email string) (c *models.User, err error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	r := "SELECT b64, email, name, password FROM users WHERE email = ?;"
	row := db.QueryRow(r, email)
	c = new(models.User)
	err = row.Scan(&c.B64, &c.Email, &c.Name, &c.Password)
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

func GetAllUsers() ([]models.ProtectedUser, error) {
	res := []models.ProtectedUser{}
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	r := "SELECT b64, name FROM users;"
	rows, err := db.Query(r)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := models.ProtectedUser{}
		err = rows.Scan(&u.B64, &u.Name)
		if err != nil {
			return nil, err
		}
		res = append(res, u)
	}

	return res, nil
}
