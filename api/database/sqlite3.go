package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"real-time-forum/api/models"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const TransactionTimeout = 3 * time.Second

type Sqlite3Store struct {
	db *sql.DB
}

func NewSqlite3Store() (*Sqlite3Store, error) {
	db, err := sql.Open("sqlite3", "api/database/database.sqlite")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		name TEXT,
		password BLOB,
		gender TEXT,
		age INTEGER,
		first_name TEXT,
		last_name TEXT,
		created DATE
	);
	
	CREATE TABLE IF NOT EXISTS posts (
		id TEXT PRIMARY KEY,
		userid TEXT REFERENCES users(id),
		categories BLOB,
		content TEXT,
		created DATE
	);`)
	if err != nil {
		return nil, err
	}

	return &Sqlite3Store{
		db: db,
	}, nil
}

func (store *Sqlite3Store) RegisterUser(req *models.RegisterRequest) (*models.User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	crypt, err := bcrypt.GenerateFromPassword([]byte(req.Password), 11)
	if err != nil {
		return nil, err
	}

	tx, err := store.db.BeginTx(req.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user := &models.User{
		ID:        id.String(),
		Email:     req.Email,
		Name:      req.Name,
		Gender:    req.Gender,
		Age:       req.Age,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Created:   time.Now().UTC(),
	}

	_, err = tx.ExecContext(req.Ctx, "INSERT INTO users VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);",
		user.ID,
		user.Email,
		user.Name,
		crypt,
		user.Gender,
		user.Age,
		user.FirstName,
		user.LastName,
		user.Created,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (store *Sqlite3Store) LogUser(req *models.LoginRequest) (*models.User, error) {
	tx, err := store.db.BeginTx(req.Ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(req.Ctx, "SELECT * FROM users WHERE email = ?;", req.Email)
	password := []byte{}
	user := new(models.User)

	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&password,
		&user.Gender,
		&user.Age,
		&user.FirstName,
		&user.LastName,
		&user.Created,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(password, []byte(req.Password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (store *Sqlite3Store) GetUsers(ctx context.Context) ([]*models.User, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, "SELECT id, email, name, gender, age, first_name , last_name, created FROM users;")
	if err != nil {
		return nil, err
	}

	users := []*models.User{}

	for rows.Next() {
		user := new(models.User)
		err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Age,
			&user.Gender,
			&user.FirstName,
			&user.LastName,
			&user.Created,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (store *Sqlite3Store) CreatePost(req *models.PostRequest) (*models.Post, error) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890+-")
	id := make([]rune, 5)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}

	tx, err := store.db.BeginTx(req.Ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	post := &models.Post{
		ID:         string(id),
		UserID:     req.UserID,
		Username:   req.Username,
		Categories: req.Categories,
		Content:    req.Content,
		Created:    time.Now().UTC(),
	}

	byteCategories, err := json.Marshal(post.Categories)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(req.Ctx, "INSERT INTO posts VALUES (?, ?, ? ,?, ?);",
		post.ID,
		post.UserID,
		byteCategories,
		post.Content,
		post.Created,
	)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (store *Sqlite3Store) GetPosts(ctx context.Context) ([]*models.Post, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx,
		`SELECT posts.id, users.id, users.name, posts.categories, posts.content, posts.created 
			FROM posts JOIN users ON posts.userid = users.id;`)
	if err != nil {
		return nil, err
	}

	posts := []*models.Post{}
	byteCategories := []byte{}

	for rows.Next() {
		post := new(models.Post)
		err = rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Username,
			&byteCategories,
			&post.Content,
			&post.Created,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(byteCategories, &post.Categories)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return posts, nil
}