package middlewares

import (
	"context"
	"net/http"
	"relif/bff/http/cookies"
	"relif/bff/services"
)

type AuthenticateByCookie struct {
	authService services.Authentication
}

func NewAuthenticateByCookie(authService services.Authentication) *AuthenticateByCookie {
	return &AuthenticateByCookie{
		authService: authService,
	}
}

func (middleware *AuthenticateByCookie) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookies.SessionCookieName)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sessionId := cookie.Value
		user, err := middleware.authService.AuthenticateSession(sessionId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)
		ctx = context.WithValue(ctx, "sessionId", sessionId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
