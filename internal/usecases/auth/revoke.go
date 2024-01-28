package usecases

import (
	"database/sql"
	"hyperzoop/internal/core/ports"

	"go.uber.org/zap"
)

type RevokeUseCase struct {
	sessionRepo ports.SessionRepository
}

// NewRevokeUseCase creates a new RevokeUseCase with the given session repository.
//
// sessionRepo: the session repository for the use case.
// *RevokeUseCase: the pointer to the newly created RevokeUseCase.
func NewRevokeUseCase(sessionRepo ports.SessionRepository) *RevokeUseCase {
	return &RevokeUseCase{sessionRepo: sessionRepo}
}

// Execute revokes a session for the given session ID and logged user.
// It returns an error.
func (r *RevokeUseCase) Execute(sessionId, loggedUser string) error {
	zap.L().Info("revoke request", zap.String("session_id", sessionId), zap.String("logged_user", loggedUser))
	session, err := r.sessionRepo.One(sessionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errSessionNotFound
		}
		zap.L().Error("error finding session", zap.Error(err))
		return err
	}
	if session.UserId != loggedUser {
		return errUnauthorized
	}
	err = r.sessionRepo.Disconnect(sessionId)
	if err != nil {
		zap.L().Error("error disconnecting session", zap.Error(err))
		return err
	}
	return nil
}
