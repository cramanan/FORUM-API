package models

import (
	"net/mail"

	"github.com/gofrs/uuid"
)

type Client struct {
	Uuid     uuid.UUID
	Email    *mail.Address
	Username string
	Password string
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
