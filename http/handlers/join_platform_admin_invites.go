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
	joinPlatformAdminInvitesUseCases "relif/platform-bff/usecases/join_platform_admin_invites"
	"relif/platform-bff/utils"
	"strconv"
)

type JoinPlatformAdminInvites struct {
	createUseCase            joinPlatformAdminInvitesUseCases.Create
	findManyPaginatedUseCase joinPlatformAdminInvitesUseCases.FindManyPaginated
	consumeByCodeUseCase     joinPlatformAdminInvitesUseCases.ConsumeByCode
}

func NewJoinPlatformAdminInvites(
	createUseCase joinPlatformAdminInvitesUseCases.Create,
	findManyPaginatedUseCase joinPlatformAdminInvitesUseCases.FindManyPaginated,
	consumeByCodeUseCase joinPlatformAdminInvitesUseCases.ConsumeByCode,
) *JoinPlatformAdminInvites {
	return &JoinPlatformAdminInvites{
		createUseCase:            createUseCase,
		findManyPaginatedUseCase: findManyPaginatedUseCase,
		consumeByCodeUseCase:     consumeByCodeUseCase,
	}
}

func (handler *JoinPlatformAdminInvites) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateJoinPlatformAdminInvite

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
	invite, err := handler.createUseCase.Execute(user, data)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInviteAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewJoinPlatformAdminInvite(invite)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinPlatformAdminInvites) FindManyPaginated(w http.ResponseWriter, r *http.Request) {
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

	count, invites, err := handler.findManyPaginatedUseCase.Execute(user, int64(offset), int64(limit))

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrForbiddenAction):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.FindMany[responses.JoinPlatformAdminInvites]{Count: count, Data: responses.NewJoinPlatformAdminInvites(invites)}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinPlatformAdminInvites) ConsumeByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	invite, err := handler.consumeByCodeUseCase.Execute(code)

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrJoinPlatformAdminInviteNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res := responses.NewJoinPlatformAdminInvite(invite)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
