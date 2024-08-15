package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"relif/platform-bff/entities"
	"relif/platform-bff/http/responses"
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
	"strconv"
)

type StorageRecords struct {
	service              services.StorageRecords
	authorizationService services.Authorization
}

func NewStorageRecords(service services.StorageRecords, authorizationService services.Authorization) *StorageRecords {
	return &StorageRecords{
		service:              service,
		authorizationService: authorizationService,
	}
}

func (handler *StorageRecords) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
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

	count, records, err := handler.service.FindManyByOrganizationID(organizationID, int64(offset), int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.StorageRecordsByLocation]{Data: responses.NewStorageRecordsByLocation(records), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *StorageRecords) FindManyByHousingID(w http.ResponseWriter, r *http.Request) {
	housingID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeAccessHousingData(housingID, user); err != nil {
		switch {
		case errors.Is(err, utils.ErrUnauthorizedAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrHousingNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	count, records, err := handler.service.FindManyByHousingID(housingID, int64(offset), int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.StorageRecordsByLocation]{Data: responses.NewStorageRecordsByLocation(records), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
