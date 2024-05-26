package models

import (
	"time"
)

type Account struct {
	Id                string    `json:"id"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	UserId            string    `json:"userId"`
	Image             string    `json:"image"`
	Type              string    `json:"type"`
	Provider          string    `json:"provider"`
	ProviderAccountId string    `json:"providerAccountId"`
	RefreshToken      string    `json:"refreshToken"`
	AccessToken       string    `json:"accessToken"`
	ExpiresAt         int64     `json:"expiresAt"`
	IdToken           string    `json:"idToken"`
	Scope             string    `json:"scope"`
	SessionState      string    `json:"sessionState"`
	TokenType         string    `json:"tokenType"`
}
