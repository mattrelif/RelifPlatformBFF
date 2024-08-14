package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"relif/platform-bff/entities"
	"relif/platform-bff/http/requests"
	"relif/platform-bff/http/responses"
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
	"strconv"
)

type JoinOrganizationInvites struct {
	service              services.JoinOrganizationInvites
	authorizationService services.Authorization
}

func NewJoinOrganizationInvites(service services.JoinOrganizationInvites, authorizationService services.Authorization) *JoinOrganizationInvites {
	return &JoinOrganizationInvites{
		service:              service,
		authorizationService: authorizationService,
	}
}

func (handler *JoinOrganizationInvites) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateJoinOrganizationInvite

	organizationID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeCreateOrganizationResource(user, organizationID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

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

	invite, err := handler.service.Create(user, req.ToEntity())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewJoinOrganizationInvite(invite)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinOrganizationInvites) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
	organizationID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeAccessPrivateOrganizationData(organizationID, user); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	offsetParam := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetParam)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limitParam := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitParam)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	count, joinRequests, err := handler.service.FindManyByOrganizationID(organizationID, int64(offset), int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.JoinOrganizationInvites]{Data: responses.NewJoinOrganizationInvites(joinRequests), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinOrganizationInvites) FindManyByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeAccessUserResource(userID, user); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	offsetParam := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetParam)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limitParam := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitParam)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	count, joinRequests, err := handler.service.FindManyByUserID(userID, int64(offset), int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.JoinOrganizationInvites]{Data: responses.NewJoinOrganizationInvites(joinRequests), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinOrganizationInvites) Accept(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeMutateJoinOrganizationInviteData(id, user); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrJoinOrganizationInviteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := handler.service.Accept(id); err != nil {
		switch {
		case errors.Is(err, utils.ErrJoinOrganizationInviteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *JoinOrganizationInvites) Reject(w http.ResponseWriter, r *http.Request) {
	var req requests.RejectJoinOrganizationInvite

	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeMutateJoinOrganizationInviteData(id, user); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrJoinOrganizationInviteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = handler.service.Reject(id, req.ToEntity()); err != nil {
		switch {
		case errors.Is(err, utils.ErrJoinOrganizationInviteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
