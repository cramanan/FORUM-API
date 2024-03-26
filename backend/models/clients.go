package models

import (
	"fmt"
	"net/mail"

	"github.com/gofrs/uuid"
)

type Client struct {
	Uuid     uuid.UUID
	Email    *mail.Address
	Username string
	Password string
}

func (c Client) String() string {
	return fmt.Sprintf("{ Uuid : %s, Email : %s, Username : %s, Password : %s}",
		c.Uuid.String(),
		c.Email.Address,
		c.Username,
		c.Password)
}

// type RawClient struct {
// 	Uuid     string
// 	Email    string
// 	Username string
// 	Password string
// }

// func (c Client) toRaw() (rc RawClient) {
// 	rc.Uuid = c.String()
// 	rc.Email = c.Email.Address
// 	rc.Username = c.Username
// 	rc.Password = c.Password
// 	return rc
// }
