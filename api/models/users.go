package models

type User struct {
	B64       string `json:"b64"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	password  string
	Gender    string `json:"gender"`
	Age       int    `json:"age"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) SetPassword(password string) {
	u.password = password
}
