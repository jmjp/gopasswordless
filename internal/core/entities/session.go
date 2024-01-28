package entities

import (
	"time"
)

type Session struct {
	Id               string    `json:"id"`
	UserId           string    `json:"user_id"`
	ValidUntil       time.Time `json:"valid_until"`
	UserAgent        *string   `json:"user_agent" bson:"user_agent"`
	Ip               *string   `json:"ip"`
	Latitude         *float64  `json:"latitude" bson:"latitude"`
	Longitude        *float64  `json:"longitude" bson:"longitude"`
	City             *string   `json:"city" bson:"city"`
	Region           *string   `json:"region" bson:"state"`
	Country          *string   `json:"country" bson:"country"`
	OrganizationName *string   `json:"organization_name" bson:"organization_name"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}

// NewSession creates a new session with the given userId, validUntil time, and userAgent pointer.
// It returns a pointer to the created Session.
func NewSession(userId string, validUntil time.Time, userAgent *string) *Session {
	return &Session{
		UserId:     userId,
		ValidUntil: validUntil,
		UserAgent:  userAgent,
	}
}

// IsExpired checks if the session is expired.
//
// No parameters.
// Returns a boolean.
func (s *Session) IsExpired() bool {
	return s.ValidUntil.Before(time.Now())
}

func (s *Session) UpdateLocation(latitude, longitude *float64, ip, city, region, country, organizationName *string) {
	s.Latitude = latitude
	s.Longitude = longitude
	s.City = city
	s.Region = region
	s.Country = country
	s.OrganizationName = organizationName
	s.Ip = ip
}
