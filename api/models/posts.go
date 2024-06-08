package models

import "time"

type Post struct {
	ID       string    `json:"id"`
	UserID   string    `json:"userID"`
	Username string    `json:"username"`
	Content  string    `json:"content"`
	Created  time.Time `json:"created"`
}

type PostRequest struct {
	Content  string `json:"content"`
	UserID   string `json:"userID"`
	Username string `json:"username"`
}
