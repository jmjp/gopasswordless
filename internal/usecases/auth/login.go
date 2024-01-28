package usecases

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"hyperzoop/internal/core/entities"
	"hyperzoop/internal/core/ports"
	"hyperzoop/internal/infra/dtos"
	"os"
	"time"

	"go.uber.org/zap"
)

type LoginUseCase struct {
	userRepository  ports.UserRepository
	magicRepository ports.MagicLinkRepository
}

// NewLoginUseCase creates a new LoginUseCase.
//
// userRepository: ports.UserRepository
// magicRepository: ports.MagicLinkRepository
// *LoginUseCase
func NewLoginUseCase(userRepository ports.UserRepository, magicRepository ports.MagicLinkRepository) *LoginUseCase {
	return &LoginUseCase{
		userRepository:  userRepository,
		magicRepository: magicRepository,
	}
}

// Execute executes the LoginUseCase with the given input and returns the LoginOutputDTO and error.
//
// input dtos.LoginInputDTO
// *dtos.LoginOutputDTO, error
func (u *LoginUseCase) Execute(input dtos.LoginInputDTO) (*dtos.LoginOutputDTO, error) {
	zap.L().Info("login request", zap.String("email", input.Email))
	user, err := u.userRepository.FindByEmail(input.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		userEntity, err := entities.NewUser(input.Email, input.Avatar, input.Username)
		if err != nil {
			return nil, err
		}
		user, err = u.userRepository.Create(userEntity)
		if err != nil {
			return nil, errCreateUser
		}
	}

	if user.Blocked {
		return nil, fmt.Errorf("user %s is blocked", user.Username)
	}

	code, fingerprint, err := generateHashedCodes()
	if err != nil {
		zap.L().Error("error generating code", zap.Error(err))
		return nil, err
	}

	if err := u.magicRepository.Create(entities.NewMagicLink(user.ID, *code, *fingerprint)); err != nil {
		zap.L().Error("error creating magic link", zap.Error(err))
	}

	return &dtos.LoginOutputDTO{
		Message:   fmt.Sprintf("magic link sent to %s", user.Email),
		Link:      fmt.Sprintf("%s/auth/verify?code=%s", os.Getenv("verify_host"), *code),
		Cookie:    fingerprint,
		ExpiresIn: time.Now().Add(time.Minute * 15),
	}, nil
}

func generateHashedCodes() (*string, *string, error) {
	code := make([]byte, 64)
	if _, err := rand.Read(code); err != nil {
		return nil, nil, err
	}
	fingerCode := make([]byte, 64)
	if _, err := rand.Read(fingerCode); err != nil {
		return nil, nil, err
	}
	tokenStr := hex.EncodeToString(code)
	fingerPrint := hex.EncodeToString(fingerCode)
	return &tokenStr, &fingerPrint, nil
}
