package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"hyperzoop/internal/core/entities"
	"hyperzoop/internal/core/ports"
	"hyperzoop/internal/infra/dtos"
	"hyperzoop/internal/infra/iplocation"
	"hyperzoop/internal/infra/token"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type AuthService struct {
	userRepository    ports.UserRepository
	magicRepository   ports.MagicLinkRepository
	sessionRepository ports.SessionRepository
}

func NewAuthService(
	userRepository ports.UserRepository,
	magicRepository ports.MagicLinkRepository,
	sessionRepository ports.SessionRepository,
) *AuthService {
	return &AuthService{
		userRepository:    userRepository,
		magicRepository:   magicRepository,
		sessionRepository: sessionRepository,
	}
}

var (
	errCreateUser               = errors.New("error creating user, please verify your email and try again")
	errInvalidCodeOrFingerprint = errors.New("verification code is invalid or fingerprint is missing")
	errNoCodeFounded            = errors.New("no magic link found, please login again")
	errUserNotFound             = errors.New("user not found")
	errSessionNotFound          = errors.New("session not found or already expired")
	errUnauthorized             = errors.New("you are not authorized to perform this action")
)

func (u *AuthService) Login(input dtos.LoginInputDTO) (*dtos.LoginOutputDTO, error) {
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

func (u *AuthService) Refresh(refresh string) (out *dtos.RefreshOutputDTO, err error) {
	zap.L().Info("refresh request", zap.String("refresh", refresh))
	session, err := u.sessionRepository.One(refresh)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errSessionNotFound
		}
		return
	}
	if session.IsExpired() {
		return nil, errSessionNotFound
	}
	user, err := u.userRepository.FindById(session.UserId)
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
			err := u.sessionRepository.Update(session)
			if err != nil {
				zap.L().Error("failed to update session", zap.Error(err))
			}
		}()
		out.RefreshToken = &session.Id
		out.ExpiresIn = &session.ValidUntil
	}
	return
}

func (u *AuthService) Revoke(sessionId, loggedUser string) error {
	zap.L().Info("revoke request", zap.String("session_id", sessionId), zap.String("logged_user", loggedUser))
	session, err := u.sessionRepository.One(sessionId)
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
	err = u.sessionRepository.Disconnect(sessionId)
	if err != nil {
		zap.L().Error("error disconnecting session", zap.Error(err))
		return err
	}
	return nil
}

func (u *AuthService) Verify(code, cookie, ip, ua string) (*dtos.VerifyOutputDTO, error) {
	zap.L().Info("verify request", zap.String("code", code), zap.String("cookie", cookie), zap.String("ip", ip), zap.String("ua", ua))
	if code == "" || len(code) < 20 || cookie == "" || len(cookie) < 20 {
		return nil, errInvalidCodeOrFingerprint
	}
	magic, err := u.magicRepository.FindValidByCode(code, cookie)
	if err != nil {
		zap.L().Error("error finding magic link", zap.Error(err))
		if err != sql.ErrNoRows {
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
		zap.L().Error("error finding user", zap.Error(err))
		if err == sql.ErrNoRows {
			return nil, errUserNotFound
		}
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

func (u *AuthService) Sessions(userID, currentToken string) ([]*dtos.SessionsOutput, error) {
	zap.L().Info("sessions request", zap.String("user_id", userID), zap.String("current_token", currentToken))
	sessions, err := u.sessionRepository.All(userID)
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

func generateAccessToken(user *entities.User) (ac string, err error) {
	ac, err = token.NewJwtAccessToken(token.UserClaims{UserId: user.ID, Email: user.Email, Blocked: user.Blocked, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(15 * time.Minute).Unix()}})
	if err != nil {
		zap.L().Error("error generating access token", zap.Error(err), zap.String("user_id", user.ID))
		return
	}
	return
}

func (u *AuthService) updateSessionGeoLocation(session *entities.Session, user *entities.User, ip string) {
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

func (u *AuthService) createSessionAndAccessToken(user *entities.User, ua string) (*entities.Session, string, error) {
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
