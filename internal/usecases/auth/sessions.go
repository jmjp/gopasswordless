package usecases

import (
	"database/sql"
	"hyperzoop/internal/core/ports"
	"hyperzoop/internal/infra/dtos"

	"go.uber.org/zap"
)

type SessionUseCase struct {
	sessionRepo ports.SessionRepository
}

// NewSessionUseCase creates a new SessionUseCase instance.
//
// It takes a sessionRepo of type ports.SessionRepository as a parameter and returns a pointer to SessionUseCase.
func NewSessionUseCase(sessionRepo ports.SessionRepository) *SessionUseCase {
	return &SessionUseCase{sessionRepo: sessionRepo}
}

// Execute fetches sessions for a user and returns the sessions with the current session marked.
//
// userID - the ID of the user for fetching sessions.
// currentToken - the current session token.
// []*dtos.SessionsOutput - a slice of sessions with the current session marked.
// error - an error, if any.
func (r *SessionUseCase) Execute(userID string, currentToken string) ([]*dtos.SessionsOutput, error) {
	zap.L().Info("sessions request", zap.String("user_id", userID), zap.String("current_token", currentToken))
	sessions, err := r.sessionRepo.All(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errSessionNotFound
		} else {
			zap.L().Error("error finding sessions", zap.Error(err))
			return nil, err
		}

	}
	var output []*dtos.SessionsOutput
	for _, session := range sessions {
		output = append(output, &dtos.SessionsOutput{
			Session: *session,
			Current: session.Id == currentToken,
		})
	}
	return output, nil
}
