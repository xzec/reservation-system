package models

import "time"

type VerificationToken struct {
	Identifier string    `json:"identifier"`
	Expires    time.Time `json:"expires"`
	Token      string    `json:"token"`
}
