package models

type User struct {
	B64       string
	Email     string
	Name      string
	Password  string
	Gender    string
	Age       int
	FirstName string
	LastName  string
}

type ProtectedUser struct {
	B64  string `json:"b64"`
	Name string `json:"name"`
}
