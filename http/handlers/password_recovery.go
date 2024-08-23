package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"relif/platform-bff/http/requests"
	passwordRecoveryUseCases "relif/platform-bff/usecases/password_recovery"
	"relif/platform-bff/utils"
)

type PasswordRecovery struct {
	requestChangeUseCase passwordRecoveryUseCases.RequestChange
	changeUseCase        passwordRecoveryUseCases.Change
}

func NewPassword(
	requestChangeUseCase passwordRecoveryUseCases.RequestChange,
	changeUseCase passwordRecoveryUseCases.Change,
) *PasswordRecovery {
	return &PasswordRecovery{
		requestChangeUseCase: requestChangeUseCase,
		changeUseCase:        changeUseCase,
	}
}

func (handler *PasswordRecovery) RequestChange(w http.ResponseWriter, r *http.Request) {
	var req requests.RequestPasswordChange

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

	if err = handler.requestChangeUseCase.Execute(req.Email); err != nil {
		switch {
		case errors.Is(err, utils.ErrMemberOfInactiveOrganization):
			http.Error(w, err.Error(), http.StatusGone)
		case errors.Is(err, utils.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, utils.ErrInactiveUser):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *PasswordRecovery) Update(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdatePassword

	code := chi.URLParam(r, "code")
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

	if err = handler.changeUseCase.Execute(code, req.NewPassword); err != nil {
		switch {
		case errors.Is(err, utils.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
