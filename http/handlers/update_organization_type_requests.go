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

type UpdateOrganizationTypeRequests struct {
	service              services.UpdateOrganizationTypeRequests
	authorizationService services.Authorization
}

func NewUpdateOrganizationTypeRequests(service services.UpdateOrganizationTypeRequests, authorizationService services.Authorization) *UpdateOrganizationTypeRequests {
	return &UpdateOrganizationTypeRequests{
		service:              service,
		authorizationService: authorizationService,
	}
}

func (handler *UpdateOrganizationTypeRequests) Create(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeCreateUpdateOrganizationTypeRequest(user); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	request, err := handler.service.Create(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewUpdateOrganizationTypeRequest(request)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *UpdateOrganizationTypeRequests) FindMany(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizePrivateActions(user); err != nil {
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

	count, reqs, err := handler.service.FindMany(int64(offset), int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.UpdateOrganizationTypeRequests]{Data: responses.NewUpdateOrganizationTypeRequests(reqs), Count: count}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *UpdateOrganizationTypeRequests) FindManyByOrganizationId(w http.ResponseWriter, r *http.Request) {
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

	count, reqs, err := handler.service.FindManyByOrganizationId(organizationId, int64(offset), int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.UpdateOrganizationTypeRequests]{Data: responses.NewUpdateOrganizationTypeRequests(reqs), Count: count}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *UpdateOrganizationTypeRequests) Accept(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizePrivateActions(user); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	id := chi.URLParam(r, "id")

	if err := handler.service.Accept(user.ID, id); err != nil {
		switch {
		case errors.Is(err, utils.ErrUpdateOrganizationTypeRequestNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *UpdateOrganizationTypeRequests) Reject(w http.ResponseWriter, r *http.Request) {
	var req requests.RejectUpdateOrganizationTypeRequest

	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizePrivateActions(user); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	id := chi.URLParam(r, "id")

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

	if err = handler.service.Reject(id, user.ID, req.ToEntity()); err != nil {
		switch {
		case errors.Is(err, utils.ErrUpdateOrganizationTypeRequestNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
