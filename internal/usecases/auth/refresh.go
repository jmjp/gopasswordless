package usecases

import (
	"database/sql"
	"hyperzoop/internal/core/ports"
	"hyperzoop/internal/infra/dtos"
	"time"

	"go.uber.org/zap"
)

type RefreshUseCase struct {
	sessionRepo ports.SessionRepository
	userRepo    ports.UserRepository
}

// NewRefreshUseCase creates a new RefreshUseCase.
//
// sessionRepo: the session repository.
// userRepository: the user repository.
// *RefreshUseCase: the new RefreshUseCase instance.
func NewRefreshUseCase(sessionRepo ports.SessionRepository, userRepository ports.UserRepository) *RefreshUseCase {
	return &RefreshUseCase{sessionRepo: sessionRepo, userRepo: userRepository}
}

// Execute executes the RefreshUseCase with the provided refresh token.
//
// refresh string
// *dtos.RefreshOutputDTO, error
func (r *RefreshUseCase) Execute(refresh string) (out *dtos.RefreshOutputDTO, err error) {
	zap.L().Info("refresh request", zap.String("refresh", refresh))
	session, err := r.sessionRepo.One(refresh)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errSessionNotFound
		}
		return
	}
	if session.IsExpired() {
		return nil, errSessionNotFound
	}
	user, err := r.userRepo.FindById(session.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errUserNotFound
		}
		zap.L().Error("error finding user", zap.Error(err))
		return
	}
	if user.Blocked {
		return nil, errUnauthorized
	}
	accessToken, err := generateAccessToken(user)
	if err != nil {
		return
	}
	out = &dtos.RefreshOutputDTO{
		AccessToken: accessToken,
		User:        *user,
	}
	if session.ValidUntil.Add(time.Hour * 12).Before(time.Now()) {
		session.ValidUntil = time.Now().Add(time.Hour * 24)
		go func() {
			err := r.sessionRepo.Update(session)
			if err != nil {
				zap.L().Error("failed to update session", zap.Error(err))
			}
		}()
		out.RefreshToken = &session.Id
		out.ExpiresIn = &session.ValidUntil
	}
	return

}
