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
	updateOrganizationTypeRequestsUseCase "relif/platform-bff/usecases/update_organization_type_requets"
	"relif/platform-bff/utils"
	"strconv"
)

type UpdateOrganizationTypeRequests struct {
	createUseCase                            updateOrganizationTypeRequestsUseCase.Create
	findManyPaginatedUseCase                 updateOrganizationTypeRequestsUseCase.FindManyPaginated
	findManyByOrganizationIDPaginatedUseCase updateOrganizationTypeRequestsUseCase.FindManyByOrganizationIDPaginated
	acceptUseCase                            updateOrganizationTypeRequestsUseCase.Accept
	rejectUseCase                            updateOrganizationTypeRequestsUseCase.Reject
}

func NewUpdateOrganizationTypeRequests(
	createUseCase updateOrganizationTypeRequestsUseCase.Create,
	findManyPaginatedUseCase updateOrganizationTypeRequestsUseCase.FindManyPaginated,
	findManyByOrganizationIDPaginatedUseCase updateOrganizationTypeRequestsUseCase.FindManyByOrganizationIDPaginated,
	acceptUseCase updateOrganizationTypeRequestsUseCase.Accept,
	rejectUseCase updateOrganizationTypeRequestsUseCase.Reject,
) *UpdateOrganizationTypeRequests {
	return &UpdateOrganizationTypeRequests{
		createUseCase:                            createUseCase,
		findManyPaginatedUseCase:                 findManyPaginatedUseCase,
		findManyByOrganizationIDPaginatedUseCase: findManyByOrganizationIDPaginatedUseCase,
		acceptUseCase:                            acceptUseCase,
		rejectUseCase:                            rejectUseCase,
	}
}

func (handler *UpdateOrganizationTypeRequests) Create(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(entities.User)

	request, err := handler.createUseCase.Execute(user)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	count, reqs, err := handler.findManyPaginatedUseCase.Execute(user, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.UpdateOrganizationTypeRequests]{Data: responses.NewUpdateOrganizationTypeRequests(reqs), Count: count}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *UpdateOrganizationTypeRequests) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
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

	count, reqs, err := handler.findManyByOrganizationIDPaginatedUseCase.Execute(user, organizationID, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.UpdateOrganizationTypeRequests]{Data: responses.NewUpdateOrganizationTypeRequests(reqs), Count: count}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *UpdateOrganizationTypeRequests) Accept(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.acceptUseCase.Execute(user, id); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
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

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = handler.rejectUseCase.Execute(user, id, req.RejectReason); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrUpdateOrganizationTypeRequestNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
