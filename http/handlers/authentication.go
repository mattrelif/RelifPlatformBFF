package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"relif/platform-bff/entities"
	"relif/platform-bff/http/requests"
	"relif/platform-bff/http/responses"
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
)

type Authentication struct {
	service services.Authentication
}

func NewAuthentication(service services.Authentication) *Authentication {
	return &Authentication{
		service: service,
	}
}

func (handler *Authentication) SignUp(w http.ResponseWriter, r *http.Request) {
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

	token, err := handler.service.SignUp(req.ToEntity())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Token", token)

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Authentication) OrganizationSignUp(w http.ResponseWriter, r *http.Request) {
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

	token, err := handler.service.OrganizationSignUp(req.ToEntity())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Token", token)

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Authentication) SignIn(w http.ResponseWriter, r *http.Request) {
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

	token, err := handler.service.SignIn(req.Email, req.Password)

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

	w.Header().Set("Token", token)

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Authentication) Me(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entities.User)

	res := responses.NewUser(user)

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Authentication) SignOut(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(entities.Session)

	if err := handler.service.SignOut(session.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
