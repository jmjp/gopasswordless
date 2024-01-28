package entities

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    *string   `json:"avatar"`
	Blocked   bool      `json:"blocked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(email string, avatar, username *string) (*User, error) {
	user := &User{
		Email:     email,
		Avatar:    avatar,
		Username:  generateUsername("hyperzoop"),
		Blocked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if username != nil {
		user.Username = *username
	}
	if err := user.isValid(); err != nil {
		return nil, err
	}
	return user, nil
}
func (user *User) isValid() error {
	usernamePattern := `^[a-zA-Z0-9_-]+$`
	emailPattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	usernameMatch, _ := regexp.MatchString(usernamePattern, user.Username)
	if !usernameMatch {
		return fmt.Errorf("invalid username, it should only contain letters, numbers, underscores, and hyphens")
	}

	emailMatch, _ := regexp.MatchString(emailPattern, user.Email)
	if !emailMatch {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func generateUsername(prefix string) string {
	rand.NewSource(time.Now().UnixNano())
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	username := make([]rune, 10)
	for i := range username {
		username[i] = runes[rand.Intn(len(runes))]
	}
	return prefix + "_" + string(username)
}
