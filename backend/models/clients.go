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
