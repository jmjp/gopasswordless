package ports

import "hyperzoop/internal/core/entities"

type SessionRepository interface {
	Create(session *entities.Session) (*entities.Session, error)
	All(userId string) ([]*entities.Session, error)
	One(id string) (*entities.Session, error)
	UpdateGeoLocation(session *entities.Session) error
	Update(session *entities.Session) error
	Disconnect(sessionId string) error
}
