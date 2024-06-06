package models

import "time"

type Post struct {
	ID       int
	UserID   string
	Username string
	Content  string
	Date     time.Time
}
