package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"relif/bff/entities"
	"relif/bff/http/requests"
	"relif/bff/http/responses"
	"relif/bff/services"
	"relif/bff/utils"
	"strconv"
)

type JoinOrganizationRequests struct {
	service              services.JoinOrganizationRequests
	authorizationService services.Authorization
}

func NewJoinOrganizationRequests(service services.JoinOrganizationRequests, authorizationService services.Authorization) *JoinOrganizationRequests {
	return &JoinOrganizationRequests{
		service:              service,
		authorizationService: authorizationService,
	}
}

func (handler *JoinOrganizationRequests) Create(w http.ResponseWriter, r *http.Request) {
	organizationId := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeCreateJoinOrganizationRequest(user); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	request, err := handler.service.Create(user.ID, organizationId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewJoinOrganizationRequest(request)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinOrganizationRequests) FindManyByOrganizationId(w http.ResponseWriter, r *http.Request) {
	organizationId := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeAccessPrivateOrganizationData(organizationId, user); err != nil {
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

	count, joinRequests, err := handler.service.FindManyByOrganizationId(organizationId, int64(offset), int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.JoinOrganizationRequests]{Data: responses.NewJoinOrganizationRequests(joinRequests), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinOrganizationRequests) Accept(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeMutateJoinOrganizationRequest(id, user); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrJoinOrganizationRequestNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := handler.service.Accept(id, user); err != nil {
		switch {
		case errors.Is(err, utils.ErrJoinOrganizationRequestNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *JoinOrganizationRequests) Reject(w http.ResponseWriter, r *http.Request) {
	var req requests.RejectJoinOrganizationRequest

	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeMutateJoinOrganizationRequest(id, user); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrJoinOrganizationRequestNotFound):
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

	if err = handler.service.Reject(id, user, req.ToEntity()); err != nil {
		switch {
		case errors.Is(err, utils.ErrJoinOrganizationRequestNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
