package models

import "time"

type Post struct {
	ID       int       `json:"id"`
	UserID   string    `json:"userid"`
	Username string    `json:"username"`
	Content  string    `json:"content"`
	Date     time.Time `json:"date"`
}
