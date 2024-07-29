package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"relif/bff/http/requests"
	"relif/bff/http/responses"
	"relif/bff/services"
	"strconv"
)

type Users struct {
	service services.Users
}

func NewUsers(service services.Users) *Users {
	return &Users{
		service: service,
	}
}

func (handler *Users) FindOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := handler.service.FindOneById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewUser(user)

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Users) FindManyByOrganizationId(w http.ResponseWriter, r *http.Request) {
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

	organizationId := chi.URLParam(r, "id")

	count, users, err := handler.service.FindManyByOrganizationId(organizationId, int64(offset), int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.Users]{Data: responses.NewUsers(users), Count: count}

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Users) UpdateOne(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateUser

	id := chi.URLParam(r, "id")

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updated, err := handler.service.FindOneAndUpdateById(id, req.ToEntity())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewUser(updated)

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Users) DeleteOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := handler.service.DeleteOneById(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
