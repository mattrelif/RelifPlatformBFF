package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"relif/platform-bff/entities"
	"relif/platform-bff/http/requests"
	"relif/platform-bff/http/responses"
	authenticationUseCases "relif/platform-bff/usecases/authentication"
	"relif/platform-bff/utils"
)

type Authentication struct {
	signUpUseCase             authenticationUseCases.SignUp
	organizationSignUpUseCase authenticationUseCases.OrganizationSignUp
	signInUseCase             authenticationUseCases.SignIn
	signOutUseCase            authenticationUseCases.SignOut
}

func NewAuthentication(
	signUpUseCase authenticationUseCases.SignUp,
	organizationSignUpUseCase authenticationUseCases.OrganizationSignUp,
	signInUseCase authenticationUseCases.SignIn,
	signOutUseCase authenticationUseCases.SignOut,
) *Authentication {
	return &Authentication{
		signUpUseCase:             signUpUseCase,
		organizationSignUpUseCase: organizationSignUpUseCase,
		signInUseCase:             signInUseCase,
		signOutUseCase:            signOutUseCase,
	}
}

func (handler *Authentication) SignUp(w http.ResponseWriter, r *http.Request) {
	var req requests.SignUp

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := req.ToEntity()
	token, err := handler.signUpUseCase.Execute(data)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrUserAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Token", token)

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Authentication) OrganizationSignUp(w http.ResponseWriter, r *http.Request) {
	var req requests.OrganizationSignUp

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := req.ToEntity()
	token, err := handler.organizationSignUpUseCase.Execute(data)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrMemberOfInactiveOrganization):
			http.Error(w, err.Error(), http.StatusGone)
		case errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, utils.ErrUserAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Token", token)

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Authentication) SignIn(w http.ResponseWriter, r *http.Request) {
	var req requests.SignIn

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := handler.signInUseCase.Execute(req.Email, req.Password)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidCredentials):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, utils.ErrMemberOfInactiveOrganization):
			http.Error(w, err.Error(), http.StatusGone)
		case errors.Is(err, utils.ErrInactiveUser):
			http.Error(w, err.Error(), http.StatusForbidden)
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

	if err := handler.signOutUseCase.Execute(session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
