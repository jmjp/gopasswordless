package dtos

import (
	"hyperzoop/internal/core/entities"
	"time"
)

type LoginInputDTO struct {
	Email    string  `json:"email"`
	Username *string `json:"username"`
	Avatar   *string `json:"avatar"`
}

type LoginOutputDTO struct {
	Message   string    `json:"message"`
	Link      string    `json:"link"`
	Cookie    *string   `json:"cookie"`
	ExpiresIn time.Time `json:"expires_in"`
}

type VerifyOutputDTO struct {
	User         *entities.User `json:"user"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    time.Time      `json:"expires_in"`
}

type RefreshOutputDTO struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken *string       `json:"refresh_token"`
	ExpiresIn    *time.Time    `json:"expires_in"`
	User         entities.User `json:"user"`
}

type SessionsOutput struct {
	entities.Session
	Current bool `json:"current"`
}
