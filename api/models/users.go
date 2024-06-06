package models

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

func (u *User) SetPassword(password []byte) {
	u.password = password
}
