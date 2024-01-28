package token

import (
	"os"

	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	UserId  string `json:"id"`
	Email   string `json:"email"`
	Blocked bool   `json:"blocked"`
	jwt.StandardClaims
}

func NewJwtAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(os.Getenv("token_secret")))
}

func ParseJwtAccessToken(accessToken string) (*UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_secret")), nil
	})
	if err != nil {
		return nil, err
	}
	return parsedAccessToken.Claims.(*UserClaims), nil
}
