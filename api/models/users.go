package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	B64       string `json:"b64"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	password  []byte
	Gender    string `json:"gender"`
	Age       int    `json:"age"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (u *User) GetPassword() []byte {
	return u.password
}

func (u *User) SetPassword(password []byte) (err error) {
	u.password, err = bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return err
	}
	return nil
}
