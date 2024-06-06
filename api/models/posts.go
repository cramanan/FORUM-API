package models

import "time"

type Post struct {
	UUID     string
	UserID   string
	Username string
	Content  string
	Date     time.Time
}
