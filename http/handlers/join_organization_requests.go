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
	joinOrganizationRequestsUseCases "relif/platform-bff/usecases/join_organization_requests"
	"relif/platform-bff/utils"
	"strconv"
)

type JoinOrganizationRequests struct {
	createUseCase                            joinOrganizationRequestsUseCases.Create
	findManyByOrganizationIDPaginatedUseCase joinOrganizationRequestsUseCases.FindManyByOrganizationIDPaginated
	findManyByUserIDPaginatedUseCase         joinOrganizationRequestsUseCases.FindManyByUserIDPaginated
	acceptUseCase                            joinOrganizationRequestsUseCases.Accept
	rejectUseCase                            joinOrganizationRequestsUseCases.Reject
}

func NewJoinOrganizationRequests(
	createUseCase joinOrganizationRequestsUseCases.Create,
	findManyByOrganizationIDPaginatedUseCase joinOrganizationRequestsUseCases.FindManyByOrganizationIDPaginated,
	findManyByUserIDPaginatedUseCase joinOrganizationRequestsUseCases.FindManyByUserIDPaginated,
	acceptUseCase joinOrganizationRequestsUseCases.Accept,
	rejectUseCase joinOrganizationRequestsUseCases.Reject,
) *JoinOrganizationRequests {
	return &JoinOrganizationRequests{
		createUseCase:                            createUseCase,
		findManyByOrganizationIDPaginatedUseCase: findManyByOrganizationIDPaginatedUseCase,
		findManyByUserIDPaginatedUseCase:         findManyByUserIDPaginatedUseCase,
		acceptUseCase:                            acceptUseCase,
		rejectUseCase:                            rejectUseCase,
	}
}

func (handler *JoinOrganizationRequests) Create(w http.ResponseWriter, r *http.Request) {
	organizationID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	request, err := handler.createUseCase.Execute(user, organizationID)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewJoinOrganizationRequest(request)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinOrganizationRequests) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
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

	count, joinRequests, err := handler.findManyByOrganizationIDPaginatedUseCase.Execute(user, organizationID, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.JoinOrganizationRequests]{Data: responses.NewJoinOrganizationRequests(joinRequests), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinOrganizationRequests) FindManyByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
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

	count, joinRequests, err := handler.findManyByUserIDPaginatedUseCase.Execute(user, userID, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	if err := handler.acceptUseCase.Execute(user, id); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
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
	var req requests.RejectJoinOrganizationRequest

	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = handler.rejectUseCase.Execute(user, id, req.RejectReason); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
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
