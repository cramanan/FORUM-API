package database

import (
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const DbFilePathEnv = ""

var Db *sql.DB

func GenerateUUID() string {
	return uuid.New().String()
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hashedPassword), nil
}
