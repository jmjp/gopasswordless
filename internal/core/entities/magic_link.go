package entities

import (
	"time"
)

type MagicLink struct {
	Id         string    `json:"id"`
	Code       string    `json:"code"`
	UserId     string    `json:"user_id"`
	Cookie     string    `json:"cookie"`
	ValidUntil time.Time `json:"valid_until"`
	Used       bool      `json:"used"`
}

// NewMagicLink creates a new MagicLink object.
//
// It takes four parameters: token (string), userId (string), cookieToken (string), and validUntil (time.Time).
//
// It returns a pointer to a MagicLink object and an error.
func NewMagicLink(userId, Code, Cookie string) *MagicLink {
	return &MagicLink{
		UserId:     userId,
		Code:       Code,
		Cookie:     Cookie,
		ValidUntil: time.Now().Add(time.Minute * 5),
		Used:       false,
	}
}

func (m *MagicLink) IsExpired() bool {
	return m.ValidUntil.Before(time.Now())
}

func (m *MagicLink) IsUsed() bool {
	return m.Used
}

func (m *MagicLink) IsValidYet() bool {
	return !m.IsExpired() && !m.IsUsed()
}

func (m *MagicLink) MarkAsUsed() *MagicLink {
	m.Used = true
	return m
}
