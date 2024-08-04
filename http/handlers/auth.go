package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"relif/bff/entities"
	"relif/bff/http/cookies"
	"relif/bff/http/requests"
	"relif/bff/http/responses"
	"relif/bff/services"
	"relif/bff/utils"
)

type Auth struct {
	service services.Authentication
}

func NewAuth(service services.Authentication) *Auth {
	return &Auth{
		service: service,
	}
}

func (handler *Auth) SignUp(w http.ResponseWriter, r *http.Request) {
	var req requests.SignUp

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := handler.service.SignUp(req.ToEntity())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := cookies.NewSessionCookie(session.SessionID, session.ExpiresAt)
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Auth) OrganizationSignUp(w http.ResponseWriter, r *http.Request) {
	var req requests.OrganizationSignUp

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := handler.service.OrganizationSignUp(req.ToEntity())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := cookies.NewSessionCookie(session.SessionID, session.ExpiresAt)
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Auth) SignIn(w http.ResponseWriter, r *http.Request) {
	var req requests.SignIn

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := handler.service.SignIn(req.Email, req.Password)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidCredentials):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, utils.ErrMemberOfInactiveOrganization):
			http.Error(w, err.Error(), http.StatusGone)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	cookie := cookies.NewSessionCookie(session.SessionID, session.ExpiresAt)
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Auth) Me(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entities.User)

	res := responses.NewUser(user)

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Auth) SignOut(w http.ResponseWriter, r *http.Request) {
	sessionId := r.Context().Value("sessionId").(string)

	if err := handler.service.SignOut(sessionId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
