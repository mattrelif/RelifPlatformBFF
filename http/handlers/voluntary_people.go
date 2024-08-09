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

type VoluntaryPeople struct {
	service              services.VoluntaryPeople
	authorizationService services.Authorization
}

func NewVoluntaryPeople(service services.VoluntaryPeople, authorizationService services.Authorization) *VoluntaryPeople {
	return &VoluntaryPeople{
		service:              service,
		authorizationService: authorizationService,
	}
}

func (handler *VoluntaryPeople) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateVoluntaryPerson

	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeCreateOrganizationResource(user); err != nil {
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

	voluntary, err := handler.service.Create(user, req.ToEntity())

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrVoluntaryPersonAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewVoluntaryPerson(voluntary)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *VoluntaryPeople) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
	organizationID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeAccessOrganizationData(organizationID, user); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	search := r.URL.Query().Get("search")

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

	count, voluntaries, err := handler.service.FindManyByOrganizationID(organizationID, search, int64(limit), int64(offset))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.VoluntaryPeople]{Data: responses.NewVoluntaryPeople(voluntaries), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *VoluntaryPeople) FindOneByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeAccessVoluntaryPersonData(id, user); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrVoluntaryPersonNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	voluntary, err := handler.service.FindOneByID(id)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrVoluntaryPersonNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewVoluntaryPerson(voluntary)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *VoluntaryPeople) Update(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateVoluntaryPerson

	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeMutateVoluntaryPersonData(id, user); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrVoluntaryPersonNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	if err = handler.service.UpdateOneByID(id, req.ToEntity()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *VoluntaryPeople) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeMutateVoluntaryPersonData(id, user); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrVoluntaryPersonNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := handler.service.InactivateOneByID(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
