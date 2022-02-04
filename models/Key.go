package models

import (
	"crypto/rsa"
	"github.com/google/uuid"
)

type Key struct {
	Revision   uuid.UUID
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}
