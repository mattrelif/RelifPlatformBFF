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
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
	"strconv"
)

type JoinPlatformInvites struct {
	service              services.JoinPlatformInvites
	authorizationService services.Authorization
}

func NewJoinPlatformInvites(service services.JoinPlatformInvites, authorizationService services.Authorization) *JoinPlatformInvites {
	return &JoinPlatformInvites{
		service:              service,
		authorizationService: authorizationService,
	}
}

func (handler *JoinPlatformInvites) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateJoinPlatformInvite

	organizationID := chi.URLParam(r, "id")
	user := r.Context().Value("user").(entities.User)

	if err := handler.authorizationService.AuthorizeCreateOrganizationResource(user, organizationID); err != nil {
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

	invite, err := handler.service.Create(req.ToEntity(), user)

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

	if err := handler.authorizationService.AuthorizeAccessPrivateOrganizationData(organizationID, user); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
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

	count, invites, err := handler.service.FindManyByOrganizationID(organizationID, int64(limit), int64(offset))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	invite, err := handler.service.ConsumeByCode(code)

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
