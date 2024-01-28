package usecases

import (
	"database/sql"
	"hyperzoop/internal/core/entities"
	"hyperzoop/internal/core/ports"
	"hyperzoop/internal/infra/dtos"
	"hyperzoop/internal/infra/iplocation"
	"hyperzoop/internal/infra/token"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type VerifyUseCase struct {
	magicRepository   ports.MagicLinkRepository
	userRepository    ports.UserRepository
	sessionRepository ports.SessionRepository
}

// NewVerifyUseCase initializes a new VerifyUseCase.
//
// magicRepository is the magic link repository.
// userRepository is the user repository.
// sessionRepository is the session repository.
// *VerifyUseCase is returned.
func NewVerifyUseCase(
	magicRepository ports.MagicLinkRepository,
	userRepository ports.UserRepository,
	sessionRepository ports.SessionRepository,
) *VerifyUseCase {
	return &VerifyUseCase{
		magicRepository:   magicRepository,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
	}
}

// Execute executes the verification process.
//
// It takes code, cookie, ip, and ua strings as parameters and returns a *VerifyOutputDTO and an error.
func (u *VerifyUseCase) Execute(code, cookie, ip, ua string) (*dtos.VerifyOutputDTO, error) {
	zap.L().Info("verify request", zap.String("code", code), zap.String("cookie", cookie), zap.String("ip", ip), zap.String("ua", ua))
	if code == "" || len(code) < 20 || cookie == "" || len(cookie) < 20 {
		return nil, errInvalidCodeOrFingerprint
	}
	magic, err := u.magicRepository.FindValidByCode(code, cookie)
	if err != nil {
		if err != sql.ErrNoRows {
			zap.L().Error("error finding magic link", zap.Error(err))
			return nil, err
		}
		return nil, errNoCodeFounded
	}
	go func() {
		err := u.magicRepository.Update(magic.MarkAsUsed())
		if err != nil {
			zap.L().Error("error invalidating magic link", zap.Error(err))
		}
	}()
	user, err := u.userRepository.FindById(magic.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errUserNotFound
		}
		zap.L().Error("error finding user", zap.Error(err))
		return nil, err
	}
	if user.Blocked {
		return nil, errUnauthorized
	}
	session, accessToken, err := u.createSessionAndAccessToken(user, ua)
	if err != nil {
		zap.L().Error("error creating session and access token", zap.Error(err))
		return nil, err
	}
	u.updateSessionGeoLocation(session, user, ip)
	return &dtos.VerifyOutputDTO{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: session.Id,
		ExpiresIn:    session.ValidUntil,
	}, nil
}

func (u *VerifyUseCase) createSessionAndAccessToken(user *entities.User, ua string) (*entities.Session, string, error) {
	var session *entities.Session
	var accessToken string
	var err error
	session, err = u.sessionRepository.Create(entities.NewSession(user.ID, time.Now().Add(time.Hour*24), &ua))
	if err != nil {
		return nil, "", err
	}
	accessToken, err = generateAccessToken(user)
	if err != nil {
		return nil, "", err
	}
	return session, accessToken, nil
}

func (u *VerifyUseCase) updateSessionGeoLocation(session *entities.Session, user *entities.User, ip string) {
	go func() {
		location, err := iplocation.GetGeoLocationByIp(ip)
		if err != nil {
			zap.L().Error("error find geolocation by ip", zap.Error(err), zap.String("ip", ip))
		}
		if location != nil {
			lat, _ := strconv.ParseFloat(location.Latitude, 64)
			long, _ := strconv.ParseFloat(location.Longitude, 64)
			session.UpdateLocation(&lat, &long, &location.Ip, location.City, location.Region, &location.Country, &location.OrganizationName)
			if err := u.sessionRepository.UpdateGeoLocation(session); err != nil {
				zap.L().Error("error updating session", zap.Error(err), zap.String("user_id", user.ID))
			}
		}
	}()
}

func generateAccessToken(user *entities.User) (ac string, err error) {
	ac, err = token.NewJwtAccessToken(token.UserClaims{UserId: user.ID, Email: user.Email, Blocked: user.Blocked, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(15 * time.Minute).Unix()}})
	if err != nil {
		zap.L().Error("error generating access token", zap.Error(err), zap.String("user_id", user.ID))
		return
	}
	return
}
