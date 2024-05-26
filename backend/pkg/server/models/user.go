package models

import "time"

type User struct {
	Id            string    `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Email         string    `json:"email"`
	EmailVerified time.Time `json:"emailVerified"`
	Name          string    `json:"name"`
	Image         string    `json:"image"`
}
