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
	organizationDataAccessRequestsUseCases "relif/platform-bff/usecases/organization_data_access_requests"
	"relif/platform-bff/utils"
	"strconv"
)

type OrganizationDataAccessRequests struct {
	createUseCase                                     organizationDataAccessRequestsUseCases.Create
	findManyByRequesterOrganizationIDPaginatedUseCase organizationDataAccessRequestsUseCases.FindManyByRequesterOrganizationIDPaginated
	findManyByTargetOrganizationIDPaginatedUseCase    organizationDataAccessRequestsUseCases.FindManyByTargetOrganizationIDPaginated
	acceptUseCase                                     organizationDataAccessRequestsUseCases.Accept
	rejectUseCase                                     organizationDataAccessRequestsUseCases.Reject
}

func NewOrganizationDataAccessRequests(
	createUseCase organizationDataAccessRequestsUseCases.Create,
	findManyByRequesterOrganizationIDPaginatedUseCase organizationDataAccessRequestsUseCases.FindManyByRequesterOrganizationIDPaginated,
	findManyByTargetOrganizationIDPaginatedUseCase organizationDataAccessRequestsUseCases.FindManyByTargetOrganizationIDPaginated,
	acceptUseCase organizationDataAccessRequestsUseCases.Accept,
	rejectUseCase organizationDataAccessRequestsUseCases.Reject,
) *OrganizationDataAccessRequests {
	return &OrganizationDataAccessRequests{
		createUseCase: createUseCase,
		findManyByRequesterOrganizationIDPaginatedUseCase: findManyByRequesterOrganizationIDPaginatedUseCase,
		findManyByTargetOrganizationIDPaginatedUseCase:    findManyByTargetOrganizationIDPaginatedUseCase,
		acceptUseCase: acceptUseCase,
		rejectUseCase: rejectUseCase,
	}
}

func (handler *OrganizationDataAccessRequests) Create(w http.ResponseWriter, r *http.Request) {
	targetOrganizationID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	request, err := handler.createUseCase.Execute(user, targetOrganizationID)

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

	res := responses.NewOrganizationDataAccessRequest(request)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *OrganizationDataAccessRequests) FindManyByRequesterOrganizationID(w http.ResponseWriter, r *http.Request) {
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

	count, reqs, err := handler.findManyByRequesterOrganizationIDPaginatedUseCase.Execute(user, organizationID, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.OrganizationDataAccessRequests]{Data: responses.NewNewOrganizationDataAccessRequests(reqs), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *OrganizationDataAccessRequests) FindManyByTargetOrganizationID(w http.ResponseWriter, r *http.Request) {
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

	count, reqs, err := handler.findManyByTargetOrganizationIDPaginatedUseCase.Execute(user, organizationID, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.OrganizationDataAccessRequests]{Data: responses.NewNewOrganizationDataAccessRequests(reqs), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *OrganizationDataAccessRequests) Accept(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.acceptUseCase.Execute(user, id); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrOrganizationDataAccessRequestNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *OrganizationDataAccessRequests) Reject(w http.ResponseWriter, r *http.Request) {
	var req requests.RejectOrganizationDataAccessRequest

	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

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

	if err = handler.rejectUseCase.Execute(user, id, req.RejectReason); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrOrganizationDataAccessRequestNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
