package services

import (
	"github.com/golang-jwt/jwt/v5"
	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
	"time"
)

type Tokens interface {
	SignToken(user entities.User, session entities.Session) (string, error)
	ParseToken(tokenString string) (string, string, error)
}

type tokensImpl struct {
	secret []byte
}

func NewTokens(secret []byte) Tokens {
	return &tokensImpl{
		secret: secret,
	}
}

func (services *tokensImpl) SignToken(user entities.User, session entities.Session) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user.ID,
		ID:        session.ID,
	})

	return token.SignedString(services.secret)
}

func (services *tokensImpl) ParseToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, utils.ErrInvalidToken
		}

		return services.secret, nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["jti"].(string), claims["sub"].(string), nil
	}

	return "", "", utils.ErrInvalidToken
}
