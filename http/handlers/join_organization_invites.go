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
	joinOrganizationInvitesUseCases "relif/platform-bff/usecases/join_organization_invites"
	"relif/platform-bff/utils"
	"strconv"
)

type JoinOrganizationInvites struct {
	createUseCase                            joinOrganizationInvitesUseCases.Create
	findManyByOrganizationIDPaginatedUseCase joinOrganizationInvitesUseCases.FindManyByOrganizationIDPaginated
	findManyByUserIDPaginatedUseCase         joinOrganizationInvitesUseCases.FindManyByUserIDPaginated
	acceptUseCase                            joinOrganizationInvitesUseCases.Accept
	rejectUseCase                            joinOrganizationInvitesUseCases.Reject
}

func NewJoinOrganizationInvites(
	createUseCase joinOrganizationInvitesUseCases.Create,
	findManyByOrganizationIDPaginatedUseCase joinOrganizationInvitesUseCases.FindManyByOrganizationIDPaginated,
	findManyByUserIDPaginatedUseCase joinOrganizationInvitesUseCases.FindManyByUserIDPaginated,
	acceptUseCase joinOrganizationInvitesUseCases.Accept,
	rejectUseCase joinOrganizationInvitesUseCases.Reject,
) *JoinOrganizationInvites {
	return &JoinOrganizationInvites{
		createUseCase:                            createUseCase,
		findManyByOrganizationIDPaginatedUseCase: findManyByOrganizationIDPaginatedUseCase,
		findManyByUserIDPaginatedUseCase:         findManyByUserIDPaginatedUseCase,
		acceptUseCase:                            acceptUseCase,
		rejectUseCase:                            rejectUseCase,
	}
}

func (handler *JoinOrganizationInvites) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateJoinOrganizationInvite

	organizationID := chi.URLParam(r, "id")
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

	data := req.ToEntity()
	invite, err := handler.createUseCase.Execute(user, organizationID, data)

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

	res := responses.NewJoinOrganizationInvite(invite)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinOrganizationInvites) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
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
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrOrganizationNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.JoinOrganizationInvites]{Data: responses.NewJoinOrganizationInvites(joinRequests), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinOrganizationInvites) FindManyByUserID(w http.ResponseWriter, r *http.Request) {
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
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.JoinOrganizationInvites]{Data: responses.NewJoinOrganizationInvites(joinRequests), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinOrganizationInvites) Accept(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.acceptUseCase.Execute(user, id); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrJoinOrganizationInviteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *JoinOrganizationInvites) Reject(w http.ResponseWriter, r *http.Request) {
	var req requests.RejectJoinOrganizationInvite

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

	if err = req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = handler.rejectUseCase.Execute(user, id, req.RejectReason); err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, utils.ErrJoinOrganizationInviteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
