package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"relif/bff/entities"
	"relif/bff/http/responses"
	"relif/bff/services"
	"relif/bff/utils"
	"strconv"
)

type JoinOrganizationRequests struct {
	service services.JoinOrganizationRequests
}

func NewJoinOrganizationRequests(service services.JoinOrganizationRequests) *JoinOrganizationRequests {
	return &JoinOrganizationRequests{
		service: service,
	}
}

func (handler *JoinOrganizationRequests) Create(w http.ResponseWriter, r *http.Request) {
	organizationId := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.service.AuthorizeCreate(user); err != nil {
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

	if err := handler.service.AuthorizeFindManyByOrganizationId(user, organizationId); err != nil {
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

	if err := handler.service.Accept(id, user); err != nil {
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

	w.WriteHeader(http.StatusNoContent)
}

func (handler *JoinOrganizationRequests) Reject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.service.Reject(id, user); err != nil {
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

	w.WriteHeader(http.StatusNoContent)
}
