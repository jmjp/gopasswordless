package ports

import (
	"hyperzoop/internal/core/entities"
	"hyperzoop/internal/infra/dtos"
)

type AuthService interface {
	Login(input dtos.LoginInputDTO) (*dtos.LoginOutputDTO, error)
	Refresh(refresh string) (out *dtos.RefreshOutputDTO, err error)
	Revoke(sessionId, loggedUser string) error
	Verify(code, cookie, ip, ua string) (*dtos.VerifyOutputDTO, error)
	Sessions(userID, currentToken string) ([]*dtos.SessionsOutput, error)
}

type SessionRepository interface {
	Create(session *entities.Session) (*entities.Session, error)
	All(userId string) ([]*entities.Session, error)
	One(id string) (*entities.Session, error)
	UpdateGeoLocation(session *entities.Session) error
	Update(session *entities.Session) error
	Disconnect(sessionId string) error
}

type MagicLinkRepository interface {
	Create(link *entities.MagicLink) error
	FindValidByCode(code, cookie string) (*entities.MagicLink, error)
	Invalidate(code string) error
	Update(link *entities.MagicLink) error
}
