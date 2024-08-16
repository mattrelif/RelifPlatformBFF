package middlewares

import (
	"context"
	"errors"
	"net/http"
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
	"strings"
)

type AuthenticateByToken struct {
	authenticationService services.Authentication
}

func NewAuthenticateByToken(authenticationService services.Authentication) *AuthenticateByToken {
	return &AuthenticateByToken{
		authenticationService: authenticationService,
	}
}

func (middleware *AuthenticateByToken) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimSpace(strings.Split(r.Header.Get("Authorization"), " ")[1])

		if token == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, session, err := middleware.authenticationService.AuthenticateToken(token)

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
		ctx = context.WithValue(ctx, "session", session)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
