package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"relif/platform-bff/entities"
	"relif/platform-bff/http/responses"
	storageRecordsUseCases "relif/platform-bff/usecases/storage_records"
	"relif/platform-bff/utils"
	"strconv"
)

type StorageRecords struct {
	findManyByOrganizationIDPaginatedUseCase storageRecordsUseCases.FindManyByOrganizationIDPaginated
	findManyByHousingIDPaginatedUseCase      storageRecordsUseCases.FindManyByHousingIDPaginated
}

func NewStorageRecords(
	findManyByOrganizationIDPaginatedUseCase storageRecordsUseCases.FindManyByOrganizationIDPaginated,
	findManyByHousingIDPaginatedUseCase storageRecordsUseCases.FindManyByHousingIDPaginated,
) *StorageRecords {
	return &StorageRecords{
		findManyByOrganizationIDPaginatedUseCase: findManyByOrganizationIDPaginatedUseCase,
		findManyByHousingIDPaginatedUseCase:      findManyByHousingIDPaginatedUseCase,
	}
}

func (handler *StorageRecords) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
	organizationID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

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

	count, records, err := handler.findManyByOrganizationIDPaginatedUseCase.Execute(user, organizationID, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	count, records, err := handler.findManyByHousingIDPaginatedUseCase.Execute(user, housingID, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrHousingNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.StorageRecordsByLocation]{Data: responses.NewStorageRecordsByLocation(records), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
