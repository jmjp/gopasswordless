package controllers

import (
	"errors"
	"fmt"
	"hyperzoop/internal/core/ports"
	"hyperzoop/internal/infra/dtos"
	"hyperzoop/internal/infra/token"
	"net/http"
	"os"
	"time"
)

type AuthenticationController struct {
	authService ports.AuthService
}

func NewAuthenticationController(authService ports.AuthService) *AuthenticationController {
	return &AuthenticationController{
		authService,
	}
}

// Login handles the authentication of a user.
//
// It takes in an http.ResponseWriter and an http.Request as parameters.
// It does not return anything.
func (c *AuthenticationController) Login(w http.ResponseWriter, r *http.Request) {
	body := new(dtos.LoginInputDTO)
	if err := RequestParseBody(r, body); err != nil {
		ResponseError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	out, err := c.authService.Login(*body)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "_fingerprint",
		Value:    *out.Cookie,
		Expires:  out.ExpiresIn,
		Path:     "/",
		Domain:   os.Getenv("app_host"),
		HttpOnly: true,
		Secure:   os.Getenv("environment") == "prod",
	})
	fmt.Println(out.Link)
	ResponseMessage(w, http.StatusOK, out.Message)
}

// Verify verifies the authentication token and fingerprint.
//
// Parameters:
// - w: http.ResponseWriter: the response writer to send the HTTP response.
// - r: *http.Request: the HTTP request received.
//
// Returns: None.
func (c *AuthenticationController) Verify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("code")
	fingerprint, err := r.Cookie("_fingerprint")
	if err != nil {
		ResponseError(w, http.StatusBadRequest, "fingerprint not found")
		return
	}
	out, err := c.authService.Verify(token, fingerprint.Value, r.RemoteAddr, r.Header.Get("User-Agent"))
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "_refresh",
		Value:    out.RefreshToken,
		Expires:  out.ExpiresIn,
		Path:     "/",
		Domain:   os.Getenv("app_host"),
		HttpOnly: true,
		Secure:   os.Getenv("environment") == "prod",
	})
	res := map[string]interface{}{"user": out.User, "access_token": out.AccessToken}
	ResponseJson(w, http.StatusOK, res)
}

// Logout logs out the user by revoking the session and deleting the refresh cookie.
//
// It takes in the following parameters:
// - w: an http.ResponseWriter used to write the response.
// - r: an *http.Request representing the incoming request.
//
// It does not return anything.
func (c *AuthenticationController) Logout(w http.ResponseWriter, r *http.Request) {
	session := r.URL.Query().Get("session")
	if session == "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "_refresh",
			Value:    "",
			Expires:  time.Unix(0, 0),
			Path:     "/",
			Domain:   os.Getenv("app_host"),
			HttpOnly: true,
			Secure:   os.Getenv("environment") == "prod",
		})
	}
	err := c.authService.Revoke(session, r.Context().Value("user").(*token.UserClaims).UserId)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Refresh handles the refreshing of authentication tokens.
//
// It takes a http.ResponseWriter and a http.Request as parameters.
// It does not return anything.
func (c *AuthenticationController) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_refresh")
	if err != nil {
		ResponseError(w, http.StatusBadRequest, errors.New("Refresh token not found").Error())
		return
	}
	out, err := c.authService.Refresh(cookie.Value)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(w, http.StatusOK, out)
}

// Sessions handles the HTTP request for user sessions.
//
// It takes in a http.ResponseWriter and a http.Request as parameters.
// It returns nothing.
func (c *AuthenticationController) Sessions(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(*token.UserClaims).UserId
	refresh, err := r.Cookie("_refresh")
	if err != nil {
		ResponseError(w, http.StatusBadRequest, errors.New("Refresh token not found").Error())
		return
	}
	out, err := c.authService.Sessions(userId, refresh.Value)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(w, http.StatusOK, out)
}
