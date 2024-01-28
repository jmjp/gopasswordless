package middlewares

import (
	"context"
	"hyperzoop/internal/infra/token"
	"net/http"
	"strings"
)

func AutheMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if len(authorization) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tks := strings.Replace(authorization, "Bearer ", "", 1)
		tks = strings.Trim(tks, " ")
		if len(tks) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		payload, err := token.ParseJwtAccessToken(tks)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", payload)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
