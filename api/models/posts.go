package models

import "time"

type Post struct {
	UserID   string
	Username string
	Content  string
	Date     time.Time
}
