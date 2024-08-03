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

type OrganizationDataAccessGrants struct {
	service services.OrganizationDataAccessGrants
}

func NewOrganizationDataAccessGrants(service services.OrganizationDataAccessGrants) *OrganizationDataAccessGrants {
	return &OrganizationDataAccessGrants{
		service: service,
	}
}

func (handler *OrganizationDataAccessGrants) FindManyByOrganizationId(w http.ResponseWriter, r *http.Request) {
	organizationId := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.service.AuthorizeFindManyByOrganizationId(user, organizationId); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

	count, grants, err := handler.service.FindManyByOrganizationId(organizationId, int64(limit), int64(offset))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.OrganizationDataAccessGrants]{Data: responses.NewOrganizationDataAccessGrants(grants), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *OrganizationDataAccessGrants) FindManyByTargetOrganizationId(w http.ResponseWriter, r *http.Request) {
	organizationId := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.service.AuthorizeFindManyByOrganizationId(user, organizationId); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

	count, grants, err := handler.service.FindManyByTargetOrganizationId(organizationId, int64(limit), int64(offset))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.OrganizationDataAccessGrants]{Data: responses.NewOrganizationDataAccessGrants(grants), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *OrganizationDataAccessGrants) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.service.AuthorizeExternalMutation(user, id); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrOrganizationDataAccessGrantNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := handler.service.DeleteOneById(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
