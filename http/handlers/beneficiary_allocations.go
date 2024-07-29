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

type BeneficiaryAllocations struct {
	service services.BeneficiaryAllocations
}

func NewBeneficiaryAllocations(service services.BeneficiaryAllocations) *BeneficiaryAllocations {
	return &BeneficiaryAllocations{
		service: service,
	}
}

func (handler *BeneficiaryAllocations) Allocate(w http.ResponseWriter, r *http.Request) {
	var req requests.AllocateBeneficiary

	user := r.Context().Value("user").(entities.User)
	beneficiaryId := chi.URLParam(r, "id")

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

	allocation, err := handler.service.Allocate(user, beneficiaryId, req.ToEntity())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewBeneficiaryAllocation(allocation)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *BeneficiaryAllocations) Reallocate(w http.ResponseWriter, r *http.Request) {
	var req requests.ReallocateBeneficiary

	user := r.Context().Value("user").(entities.User)
	beneficiaryId := chi.URLParam(r, "id")

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

	allocation, err := handler.service.Reallocate(user, beneficiaryId, req.ToEntity())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.NewBeneficiaryAllocation(allocation)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *BeneficiaryAllocations) FindManyByBeneficiaryId(w http.ResponseWriter, r *http.Request) {
	beneficiaryId := chi.URLParam(r, "id")

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

	count, allocations, err := handler.service.FindManyByBeneficiaryId(beneficiaryId, int64(offset), int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.BeneficiaryAllocations]{Data: responses.NewBeneficiaryAllocations(allocations), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *BeneficiaryAllocations) FindManyByHousingId(w http.ResponseWriter, r *http.Request) {
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

	count, allocations, err := handler.service.FindManyByHousingId(housingId, int64(offset), int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.BeneficiaryAllocations]{Data: responses.NewBeneficiaryAllocations(allocations), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *BeneficiaryAllocations) FindManyByRoomId(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "id")

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

	count, allocations, err := handler.service.FindManyByRoomId(roomId, int64(offset), int64(limit))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := responses.FindMany[responses.BeneficiaryAllocations]{Data: responses.NewBeneficiaryAllocations(allocations), Count: count}

	if err = json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
