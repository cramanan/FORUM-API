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
	Gender   string
}
