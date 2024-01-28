package usecases

import "errors"

var (
	errCreateUser               = errors.New("error creating user, please verify your email and try again")
	errInvalidCodeOrFingerprint = errors.New("verification code is invalid or fingerprint is missing")
	errNoCodeFounded            = errors.New("no magic link found, please login again")
	errUserNotFound             = errors.New("user not found")
	errSessionNotFound          = errors.New("session not found or already expired")
	errUnauthorized             = errors.New("you are not authorized to perform this action")
)
