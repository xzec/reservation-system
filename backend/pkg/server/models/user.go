package models

import "time"

type User struct {
	Id            *string    `json:"id"`
	CreatedAt     *time.Time `json:"-"`
	UpdatedAt     *time.Time `json:"-"`
	Email         *string    `json:"email"`
	EmailVerified *time.Time `json:"emailVerified"`
	Name          *string    `json:"name"`
	Image         *string    `json:"image"`
}
