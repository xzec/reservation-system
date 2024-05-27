package models

import "time"

type Session struct {
	Id           *string    `json:"id"`
	CreatedAt    *time.Time `json:"-"`
	UpdatedAt    *time.Time `json:"-"`
	UserId       *string    `json:"userId"`
	Expires      *time.Time `json:"expires"`
	SessionToken *string    `json:"sessionToken"`
}
