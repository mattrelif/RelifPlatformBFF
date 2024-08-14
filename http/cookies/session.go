package cookies

import (
	"net/http"
	"time"
)

const SessionCookieName = "Session"

func NewSessionCookie(sessionID string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/api/v1",
		HttpOnly: false,
		Secure:   true,
		Expires:  expiresAt,
		SameSite: http.SameSiteNoneMode,
	}
}
