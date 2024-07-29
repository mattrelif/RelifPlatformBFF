package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"relif/bff/entities"
	"relif/bff/http/requests"
	"relif/bff/http/responses"
	"relif/bff/services"
	"strconv"
)

type JoinPlatformInvites struct {
	service services.JoinPlatformInvites
}

func NewJoinPlatformInvites(service services.JoinPlatformInvites) *JoinPlatformInvites {
	return &JoinPlatformInvites{
		service: service,
	}
}

func (handler *JoinPlatformInvites) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateJoinPlatformInvite

	user := r.Context().Value("user").(entities.User)

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewJoinPlatformInvite(invite)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *JoinPlatformInvites) FindManyByOrganizationId(w http.ResponseWriter, r *http.Request) {
	organizationId := chi.URLParam(r, "id")

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

	count, invites, err := handler.service.FindManyByOrganizationId(organizationId, int64(limit), int64(offset))

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
