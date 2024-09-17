package middlewares

import (
	"context"
	"errors"
	"net/http"
	authenticationUseCases "relif/platform-bff/usecases/authentication"
	"relif/platform-bff/utils"
	"strings"
)

type AuthenticateByToken struct {
	authenticateTokenUseCase authenticationUseCases.AuthenticateToken
}

func NewAuthenticateByToken(authenticateTokenUseCase authenticationUseCases.AuthenticateToken) *AuthenticateByToken {
	return &AuthenticateByToken{
		authenticateTokenUseCase: authenticateTokenUseCase,
	}
}

func (middleware *AuthenticateByToken) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		splitHeader := strings.Split(header, " ")

		if len(splitHeader) != 2 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		token := splitHeader[1]

		if token == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		user, session, err := middleware.authenticateTokenUseCase.Execute(token)

		if err != nil {
			switch {
			case errors.Is(err, utils.ErrMemberOfInactiveOrganization):
				http.Error(w, err.Error(), http.StatusGone)
			case errors.Is(err, utils.ErrUserNotFound):
				http.Error(w, err.Error(), http.StatusNotFound)
			case errors.Is(err, utils.ErrInactiveUser):
				http.Error(w, err.Error(), http.StatusForbidden)
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
