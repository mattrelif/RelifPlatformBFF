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
	joinPlatformInvitesUseCases "relif/platform-bff/usecases/join_platform_invites"
	"relif/platform-bff/utils"
	"strconv"
)

type JoinPlatformInvites struct {
	createUseCase                            joinPlatformInvitesUseCases.Create
	findManyByOrganizationIDPaginatedUseCase joinPlatformInvitesUseCases.FindManyByOrganizationIDPaginated
	consumeByCodeUseCase                     joinPlatformInvitesUseCases.ConsumeByCode
}

func NewJoinPlatformInvites(
	createUseCase joinPlatformInvitesUseCases.Create,
	findManyByOrganizationIDPaginatedUseCase joinPlatformInvitesUseCases.FindManyByOrganizationIDPaginated,
	consumeByCodeUseCase joinPlatformInvitesUseCases.ConsumeByCode,
) *JoinPlatformInvites {
	return &JoinPlatformInvites{
		createUseCase:                            createUseCase,
		findManyByOrganizationIDPaginatedUseCase: findManyByOrganizationIDPaginatedUseCase,
		consumeByCodeUseCase:                     consumeByCodeUseCase,
	}
}

func (handler *JoinPlatformInvites) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateJoinPlatformInvite

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
		case errors.Is(err, utils.ErrUserAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewJoinPlatformInvite(invite)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinPlatformInvites) FindManyByOrganizationID(w http.ResponseWriter, r *http.Request) {
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

	count, invites, err := handler.findManyByOrganizationIDPaginatedUseCase.Execute(user, organizationID, int64(limit), int64(offset))

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

	res := responses.FindMany[responses.JoinPlatformInvites]{Data: responses.NewJoinPlatformInvites(invites), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinPlatformInvites) Consume(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	invite, err := handler.consumeByCodeUseCase.Execute(code)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewJoinPlatformInvite(invite)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
