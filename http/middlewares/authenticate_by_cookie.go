package middlewares

import (
	"context"
	"errors"
	"net/http"
	"relif/platform-bff/http/cookies"
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
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

		sessionID := cookie.Value
		user, err := middleware.authService.AuthenticateSession(sessionID)

		if err != nil {
			switch {
			case errors.Is(err, utils.ErrMemberOfInactiveOrganization):
				http.Error(w, err.Error(), http.StatusGone)
			case errors.Is(err, utils.ErrUserNotFound):
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)
		ctx = context.WithValue(ctx, "sessionID", sessionID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
