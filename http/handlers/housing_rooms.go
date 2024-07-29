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

type HousingRooms struct {
	service services.HousingRooms
}

func NewHousingRooms(service services.HousingRooms) *HousingRooms {
	return &HousingRooms{
		service: service,
	}
}

func (handler *HousingRooms) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateManyHousingRooms

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

	rooms, err := handler.service.CreateMany(req.ToEntity())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.HousingRooms]{Data: responses.NewHousingRooms(rooms), Count: int64(len(rooms))}

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *HousingRooms) FindOneById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	room, err := handler.service.FindOneById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewHousingRoom(room)

	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *HousingRooms) FindManyByHousingId(w http.ResponseWriter, r *http.Request) {
	housingId := chi.URLParam(r, "id")

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

	count, rooms, err := handler.service.FindManyByHousingId(housingId, int64(limit), int64(offset))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.HousingRooms]{Data: responses.NewHousingRooms(rooms), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *HousingRooms) Update(w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateHousingRoom

	id := chi.URLParam(r, "id")

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

	updated, err := handler.service.FindOneAndUpdateById(id, req.ToEntity())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewHousingRoom(updated)

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *HousingRooms) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := handler.service.DeleteOneById(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
