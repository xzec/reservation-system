package models

import "time"

type Session struct {
	Id           string    `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	UserId       string    `json:"userId"`
	Expires      time.Time `json:"expires"`
	SessionToken string    `json:"sessionToken"`
}
