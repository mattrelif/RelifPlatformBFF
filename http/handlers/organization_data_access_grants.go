package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"relif/platform-bff/entities"
	"relif/platform-bff/http/responses"
	organizationDataAccessGrantsUseCases "relif/platform-bff/usecases/organization_data_access_grants"
	"relif/platform-bff/utils"
	"strconv"
)

type OrganizationDataAccessGrants struct {
	findManyByOrganizationIDPaginatedUseCase       organizationDataAccessGrantsUseCases.FindManyByOrganizationIDPaginated
	findManyByTargetOrganizationIDPaginatedUseCase organizationDataAccessGrantsUseCases.FindManyByTargetOrganizationIDPaginated
	deleteOneByIDUseCase                           organizationDataAccessGrantsUseCases.DeleteOneByID
}

func NewOrganizationDataAccessGrants(
	findManyByOrganizationIDPaginatedUseCase organizationDataAccessGrantsUseCases.FindManyByOrganizationIDPaginated,
	findManyByTargetOrganizationIDPaginatedUseCase organizationDataAccessGrantsUseCases.FindManyByTargetOrganizationIDPaginated,
	deleteOneByIDUseCase organizationDataAccessGrantsUseCases.DeleteOneByID,
) *OrganizationDataAccessGrants {
	return &OrganizationDataAccessGrants{
		findManyByOrganizationIDPaginatedUseCase:       findManyByOrganizationIDPaginatedUseCase,
		findManyByTargetOrganizationIDPaginatedUseCase: findManyByTargetOrganizationIDPaginatedUseCase,
		deleteOneByIDUseCase:                           deleteOneByIDUseCase,
	}
}

func (handler *OrganizationDataAccessGrants) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
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

	count, grants, err := handler.findManyByOrganizationIDPaginatedUseCase.Execute(user, organizationID, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.OrganizationDataAccessGrants]{Data: responses.NewOrganizationDataAccessGrants(grants), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *OrganizationDataAccessGrants) FindManyByTargetOrganizationID(w http.ResponseWriter, r *http.Request) {
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

	count, grants, err := handler.findManyByTargetOrganizationIDPaginatedUseCase.Execute(user, organizationID, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	if err := handler.deleteOneByIDUseCase.Execute(user, id); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrOrganizationDataAccessGrantNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
